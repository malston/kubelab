package ioutil

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestHTTPDownloader_Download(t *testing.T) {
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
				path: fmt.Sprintf("%s/mkcert", os.TempDir()),
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
			if err := d.Download(
				tt.args.path,
				tt.args.url,
				0o755,
			); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
