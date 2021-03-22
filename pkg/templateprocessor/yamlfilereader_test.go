// Copyright Contributors to the Open Cluster Management project

package templateprocessor

import (
	"reflect"
	"testing"
)

func TestYamlFileReader_Asset(t *testing.T) {
	type fields struct {
		rootDirectory string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				rootDirectory: "../../test/unit/resources/yamlfilereader",
			},
			args: args{
				name: "filereader.yaml",
			},
			want: []byte(`# Copyright Contributors to the Open Cluster Management project

apiVersion: fake/v1
kind: Fake
metadata:
  name: {{ .Values }}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YamlFileReader{
				path: tt.fields.rootDirectory,
			}
			got, err := y.Asset(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("YamlFileReader.Asset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlFileReader.Asset() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestYamlFileReader_AssetNames(t *testing.T) {
	type fields struct {
		rootDirectory string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				rootDirectory: "../../test/unit/resources/yamlfilereader",
			},
			want:    []string{"filereader.yaml"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &YamlFileReader{
				path: tt.fields.rootDirectory,
			}
			got, err := r.AssetNames()
			if (err != nil) != tt.wantErr {
				t.Errorf("YamlFileReader.Asset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YamlFileReader.AssetNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
