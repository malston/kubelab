package ioutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPDownloader_Download(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	expectedFilePath := filepath.Join(tmpDir)

	type args struct {
		path string
		url  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "downloads latest mkcert",
			args: args{
				path: expectedFilePath,
				url: fmt.Sprintf(
					"https://dl.filippo.io/mkcert/latest?for=%s/%s",
					runtime.GOOS,
					runtime.GOARCH,
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := HTTPDownloader{}
			if err := d.Download(tt.args.path,
				tt.args.url, 0o755); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
