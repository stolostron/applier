// Copyright Contributors to the Open Cluster Management project

package cmd

import (
	"reflect"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Test_newOptions(t *testing.T) {
	streams := genericclioptions.IOStreams{}
	type args struct {
		streams genericclioptions.IOStreams
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		{
			name: "success",
			args: args{
				streams: streams,
			},
			want: &Options{
				ConfigFlags: genericclioptions.NewConfigFlags(true),

				IOStreams: streams,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newOptions(tt.args.streams); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
