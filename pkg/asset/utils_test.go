// Copyright Red Hat

package asset

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMemFS_ExtractAssets(t *testing.T) {
	type fields struct {
		files []string
		data  map[string][]byte
	}
	type args struct {
		prefix     string
		dir        string
		excluded   []string
		headerFile string
	}
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "2 files no execluded",
			fields: fields{
				files: []string{"file1", "file2"},
				data: map[string][]byte{
					"file1": []byte("file1content"),
					"file2": []byte("file2content"),
				},
			},
			args: args{
				prefix:   "",
				dir:      dir,
				excluded: []string{"file1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemFS{
				files: tt.fields.files,
				data:  tt.fields.data,
			}
			if err := ExtractAssets(r, tt.args.prefix, tt.args.dir, tt.args.excluded, tt.args.headerFile); (err != nil) != tt.wantErr {
				t.Errorf("MemFS.ExtractAssets() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, err := ioutil.ReadFile(filepath.Join(dir, "file2"))
			if err != nil {
				t.Error(err)
			}
			if string(b) != "file2content" {
				t.Errorf("expect %s got %s", "file2content", string(b))
			}
		})
	}
}
