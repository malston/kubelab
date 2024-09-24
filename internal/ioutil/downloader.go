package ioutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type HTTPDownloader struct{}

func NewDownloader() HTTPDownloader {
	return HTTPDownloader{}
}

func (d HTTPDownloader) Download(path, url string, mode os.FileMode) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := download(path, url)
		if err != nil {
			return fmt.Errorf("failed to download mkcert to %s, %v", path, err)
		}
	}

	err = os.Chmod(path, mode)
	if err != nil {
		return err
	}

	return nil
}

func download(path, url string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %s, failed to download file", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
