// Copyright Red Hat
package apply

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stolostron/applier/pkg/cmd/apply/common"
)

func TestOptions_Complete(t *testing.T) {
	type fields struct {
		options common.Options
	}
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "read value file succees",
			fields: fields{
				options: common.Options{
					ValuesPath: "../../../test/unit/resources/scenario/values.yaml",
				},
			},
			wantErr: false,
		},
		{
			name: "read value file not found",
			fields: fields{
				options: common.Options{
					ValuesPath: "file_not_found.yaml",
				},
			},
			wantErr: true,
		},
		{
			name:    "read value stdin",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				options: tt.fields.options,
			}
			var fileIn *os.File
			var err error
			if len(o.options.ValuesPath) == 0 {
				fileIn, err = ioutil.TempFile("", "stdin")
				if err != nil {
					t.Error(err)
				}
				defer os.Remove(fileIn.Name())
				err = ioutil.WriteFile(fileIn.Name(), []byte("ServiceAccount: my-sa\n"), 0600)
				if err != nil {
					t.Error(err)
				}
				os.Stdin = fileIn
			}
			if err := o.Complete(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Options.Complete() error = %v, wantErr %v", err, tt.wantErr)
			}
			switch tt.name {
			case "read value file succees":
				iSimple, ok := o.options.Values["Simple"]
				if !ok {
					t.Error("'Simple' not found in value")
				}
				simple := iSimple.(map[string]interface{})
				iSA, ok := simple["ServiceAccount"]
				if !ok {
					t.Error("'ServiceAccount' not found in value")
				}
				sa := iSA.(string)
				if sa != "my-sa" {
					t.Errorf("'Expected 'my-sa' got %s", sa)
				}
			case "read value stdin":
				iSA, ok := o.options.Values["ServiceAccount"]
				if !ok {
					t.Error("'ServiceAccount' not found in value")
				}
				sa := iSA.(string)
				if sa != "my-sa" {
					t.Errorf("'Expected 'my-sa' got %s", sa)
				}
			}
		})
	}
}

func TestOptions_Validate(t *testing.T) {
	type fields struct {
		options common.Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "directory succees",
			fields: fields{
				options: common.Options{
					Header: "../../../test/unit/resources/scenario/musttemplateasset/header.txt",
					Paths:  []string{"../../../test/unit/resources/scenario/musttemplateasset"},
				},
			},
			wantErr: false,
		},
		{
			name: "directory failed",
			fields: fields{
				options: common.Options{
					Header: "../../../test/unit/resources/scenario/musttemplateasset/header.txt",
					Paths:  []string{"wrong_dir"},
				},
			},
			wantErr: true,
		},
		{
			name:    "empty failed",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				options: tt.fields.options,
			}
			if err := o.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Options.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
