package cert

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/malston/kubelab/internal/ioutil"
	"github.com/stretchr/testify/require"
)

type stubExecutor struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func (s *stubExecutor) Execute(cmd string, args ...string) (string, string, error) {
	c := exec.Command(cmd, args...)
	c.Stdout = &s.Stdout
	c.Stderr = &s.Stderr

	err := c.Run()
	if err != nil {
		return "", "", err
	}

	if len(s.Stderr.String()) > 0 {
		return "", "", err
	}

	return s.Stdout.String(), "", nil
}

func TestMakeCert(t *testing.T) {
	domains := []string{"example.com"}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("something went wrong getting home path: %s", err)
	}

	configDir := fmt.Sprintf("%s/.kubelab", homeDir)
	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			log.Printf("%s directory already exists, continuing", configDir)
		}
	}

	certDir := fmt.Sprintf("%s/%s/ssl", configDir, "test-cluster")
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		err := os.MkdirAll(certDir, os.ModePerm)
		if err != nil {
			log.Printf("%s directory already exists, continuing", certDir)
		}
	}

	certFileName := certDir + "/" + domains[0] + ".crt"
	keyFileName := certDir + "/" + domains[0] + ".key"

	type args struct {
		installDir   string
		certFileName string
		keyFileName  string
		domains      []string
	}
	require.NoError(t, err)
	installer := &MkCertInstaller{
		&stubExecutor{},
		ioutil.NewDownloader(),
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "mkcert creates valid cert",
			args: args{
				installDir:   configDir,
				certFileName: certFileName,
				keyFileName:  keyFileName,
				domains:      domains,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := installer.MkCert(tt.args.installDir, tt.args.certFileName, tt.args.keyFileName, tt.args.domains...); (err != nil) != tt.wantErr {
				t.Errorf("MkCert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMkCertInstaller_Download(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	path := filepath.Join(tmpDir)

	type fields struct {
		CommandExecutor CommandExecutor
		Downloader      Downloader
	}
	type args struct {
		path string
		url  string
		mode os.FileMode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "downloads latest mkcert",
			fields: fields{
				Downloader: ioutil.NewDownloader(),
			},
			args: args{
				path: fmt.Sprintf(
					"%s/mkcert",
					path,
				),
				url: fmt.Sprintf(
					"https://dl.filippo.io/mkcert/latest?for=%s/%s",
					runtime.GOOS,
					runtime.GOARCH,
				),
				mode: 0o755,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MkCertInstaller{
				CommandExecutor: tt.fields.CommandExecutor,
				Downloader:      tt.fields.Downloader,
			}
			if err := m.Download(tt.args.path, tt.args.url, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
