package cert

import (
	"fmt"
	"os"
	"runtime"

	internalos "github.com/malston/kubelab/internal/os"
)

type MkCertInstaller struct {
	mkCert string
	Downloader
	CommandExecutor
}

type MkCertClient struct {
	Executable string
	CommandExecutor
}

type CommandExecutor interface {
	Execute(cmd string, args ...string) (string, string, error)
}

type Downloader interface {
	Download(path, url string, mode os.FileMode) error
}

func NewMkCertInstaller(dl Downloader) *MkCertInstaller {
	return &MkCertInstaller{
		Downloader:      dl,
		CommandExecutor: &internalos.ShellExecutor{},
	}
}

// Install installs mkcert to the given path
func (m MkCertInstaller) Install(path string) (*MkCertClient, error) {
	executable := fmt.Sprintf("%s/mkcert", path)
	err := m.Download(
		executable,
		fmt.Sprintf("https://dl.filippo.io/mkcert/latest?for=%s/%s",
			runtime.GOOS,
			runtime.GOARCH),
		0o755)
	if err != nil {
		return nil, fmt.Errorf("failed to download mkcert to %w", err)
	}

	return &MkCertClient{Executable: executable, CommandExecutor: m.CommandExecutor}, nil
}

// Download mkcert to the given path
func (m MkCertInstaller) Download(path, url string, mode os.FileMode) error {
	err := m.Downloader.Download(path, url, mode)
	if err != nil {
		return fmt.Errorf("failed to download mkcert to %s from %s, %v", path, url, err)
	}
	return nil
}

// MkCert invokes mkcert to create a new certificate valid for a given list of domain names
func (c MkCertClient) MkCert(certFileName string, keyFileName string, domains []string) error {
	_, _, err := c.Execute(
		c.Executable,
		append([]string{
			"-cert-file",
			certFileName,
			"-key-file",
			keyFileName,
		}, domains...)...,
	)
	return err
}
