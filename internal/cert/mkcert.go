package cert

import (
	"fmt"
	"os"
	"runtime"

	"github.com/malston/kubelab/internal/ioutil"
	internalos "github.com/malston/kubelab/internal/os"
)

type MkCertInstaller struct {
	CommandExecutor
	Downloader
}

type CommandExecutor interface {
	Execute(cmd string, args ...string) (string, string, error)
}

type Downloader interface {
	Download(path, url string, mode os.FileMode) error
}

func NewMkCertInstaller() *MkCertInstaller {
	return &MkCertInstaller{
		CommandExecutor: &internalos.ShellExecutor{},
		Downloader:      ioutil.NewDownloader(),
	}
}

func (m MkCertInstaller) Download(path, url string, mode os.FileMode) error {
	err := m.Downloader.Download(path, url, mode)
	if err != nil {
		return fmt.Errorf("failed to download mkcert to %s from %s, %v", path, url, err)
	}
	return nil
}

// MkCert invokes mkcert to create a new certificate valid for a given list of domain names
func (m MkCertInstaller) MkCert(installDir string, certFileName string, keyFileName string, domains ...string) error {
	mkCert := fmt.Sprintf("%s/mkcert", installDir)
	_, err := os.Stat(mkCert)
	if os.IsNotExist(err) {
		dlErr := m.Download(mkCert, fmt.Sprintf("https://dl.filippo.io/mkcert/latest?for=%s/%s", runtime.GOOS, runtime.GOARCH), 0o755)
		if dlErr != nil {
			return fmt.Errorf("failed to download mkcert to %s, %v", installDir, dlErr)
		}
		return m.MkCert(installDir, certFileName, keyFileName, domains...)
	}
	if err != nil {
		return fmt.Errorf("failed to stat mkcert: %w", err)
	}
	_, _, err = m.Execute(
		mkCert,
		append([]string{
			"-cert-file",
			certFileName,
			"-key-file",
			keyFileName,
		}, domains...)...,
	)

	return err
}
