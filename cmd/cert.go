package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/malston/kubelab/internal/cert"
	"github.com/spf13/cobra"
)

const (
	clusterName = "tkg-workload"
)

type mkCertOptions struct {
	domains []string
}

func defaultOptions() *mkCertOptions {
	return &mkCertOptions{}
}

func newMkCertCmd() *cobra.Command {
	o := defaultOptions()

	cmd := &cobra.Command{
		Use:          "mkcert",
		Short:        "create an SSL certificate",
		Long:         "create a new SSL certificate valid for a given set of domains",
		SilenceUsage: true,
		RunE:         o.run,
	}
	cmd.Flags().StringSliceVar(&o.domains, "domains", o.domains, "list of domains (required)")
	_ = cmd.MarkFlagRequired("domains")

	return cmd
}

func (o *mkCertOptions) run(cmd *cobra.Command, _ []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("something went wrong getting home path: %s", err)
	}

	configPath := fmt.Sprintf("%s/.kubelab", homeDir)
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			log.Printf("%s directory already exists, continuing", configPath)
		}
	}

	certPath := fmt.Sprintf("%s/%s/ssl", configPath, clusterName)
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		err := os.MkdirAll(certPath, os.ModePerm)
		if err != nil {
			log.Printf("%s directory already exists, continuing", certPath)
		}
	}

	certFileName := certPath + "/" + o.domains[0] + ".crt"
	keyFileName := certPath + "/" + o.domains[0] + ".key"

	installer := cert.NewMkCertInstaller()
	err = installer.MkCert(
		configPath,
		certFileName,
		keyFileName,
		o.domains...,
	)
	if err != nil {
		return err
	}

	if _, err = os.ReadFile(certFileName); err != nil {
		return fmt.Errorf("error reading %s file, %w", certFileName, err)
	}

	if _, err = os.ReadFile(keyFileName); err != nil {
		return fmt.Errorf("error reading %s file, %w", keyFileName, err)
	}

	return err
}
