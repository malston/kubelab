package cert

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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

func TestMkCertInstaller_Download(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	path := filepath.Join(tmpDir)

	type fields struct {
		Downloader Downloader
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
				Downloader: tt.fields.Downloader,
			}
			if err := m.Download(tt.args.path, tt.args.url, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMkCertInstaller_Install(t *testing.T) {
	type fields struct {
		mkCert          string
		Downloader      Downloader
		CommandExecutor CommandExecutor
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *MkCertClient
		wantErr bool
	}{
		{
			name: "installs mkcert",
			fields: fields{
				mkCert:          fmt.Sprintf("%s/mkcert", os.TempDir()),
				Downloader:      ioutil.NewDownloader(),
				CommandExecutor: &stubExecutor{},
			},
			args: args{path: os.TempDir()},
			want: &MkCertClient{
				Executable:      fmt.Sprintf("%s/mkcert", os.TempDir()),
				CommandExecutor: &stubExecutor{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MkCertInstaller{
				mkCert:          tt.fields.mkCert,
				Downloader:      tt.fields.Downloader,
				CommandExecutor: tt.fields.CommandExecutor,
			}
			got, err := m.Install(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Install() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMkCertClient_MkCert(t *testing.T) {
	type fields struct {
		Executable      string
		CommandExecutor CommandExecutor
	}
	type args struct {
		certFileName string
		keyFileName  string
		domains      []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "creates valid cert",
			fields: fields{
				Executable:      fmt.Sprintf("%s/mkcert", os.TempDir()),
				CommandExecutor: &stubExecutor{},
			},
			args: args{
				certFileName: fmt.Sprintf("%s/cert.pem", os.TempDir()),
				keyFileName:  fmt.Sprintf("%s/key.pem", os.TempDir()),
				domains:      []string{"example.com"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MkCertClient{
				Executable:      tt.fields.Executable,
				CommandExecutor: tt.fields.CommandExecutor,
			}
			if err := c.MkCert(tt.args.certFileName, tt.args.keyFileName, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("MkCert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
