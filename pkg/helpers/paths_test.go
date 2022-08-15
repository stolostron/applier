// Copyright Red Hat

package helpers

import (
	"testing"

	"github.com/stolostron/applier/pkg/asset"
	"github.com/stolostron/applier/test/unit/resources/scenario"
)

func TestHasMultipleAssets(t *testing.T) {
	type args struct {
		reader asset.ScenarioReader
		path   string
	}
	reader := scenario.GetScenarioResourcesReader()
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Split multicontent file",
			args: args{
				reader: reader,
				path:   "multicontent/file1.yaml",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Split single content",
			args: args{
				reader: reader,
				path:   "multicontent/clusterrole.yaml",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Split not existing file",
			args: args{
				reader: reader,
				path:   "multicontent/not-exists.yaml",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HasMultipleAssets(tt.args.reader, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasMultipleAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HasMultipleAssets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitFiles(t *testing.T) {
	type args struct {
		reader asset.ScenarioReader
		paths  []string
		Header string
	}
	reader := scenario.GetScenarioResourcesReader()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Split multicontent file",
			args: args{
				reader: reader,
				paths:  []string{"multicontent/file1.yaml"},
			},
			wantErr: false,
		},
		{
			name: "Split single content",
			args: args{
				reader: reader,
				paths:  []string{"multicontent/clusterrole.yaml"},
			},
			wantErr: false,
		},
		{
			name: "Split not existing file",
			args: args{
				reader: reader,
				paths:  []string{"multicontent/not-exists.yaml"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitFiles(tt.args.reader, tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			switch tt.name {
			case "Split multicontent file":
				if got == nil {
					t.Error("memFS is nil")
				}
				names, err := got.AssetNames(tt.args.paths, []string{}, tt.args.Header)
				if (err != nil) != tt.wantErr {
					t.Errorf("HasMultipleAssets() error = %v, wantErr %v", err, tt.wantErr)
				}
				if len(names) != 2 {
					t.Errorf("expect nb %d got %d", 2, len(names))
				}
			case "Split single content":
				if got == nil {
					t.Error("memFS is nil")
				}
				names, err := got.AssetNames(tt.args.paths, []string{}, tt.args.Header)
				if (err != nil) != tt.wantErr {
					t.Errorf("HasMultipleAssets() error = %v, wantErr %v", err, tt.wantErr)
				}
				if len(names) != 1 {
					t.Errorf("expect nb %d got %d", 2, len(names))
				}
			case "Split not existing file":
				if got != nil {
					t.Error("memFS is mpt nil")
				}
				if (err != nil) != tt.wantErr {
					t.Errorf("HasMultipleAssets() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
