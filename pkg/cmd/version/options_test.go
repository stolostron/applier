// Copyright Contributors to the Open Cluster Management project
package version

import (
	"reflect"
	"testing"

	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Test_newOptions(t *testing.T) {
	applierFlags := genericclioptionsapplier.NewApplierFlags(nil)
	type args struct {
		applierFlags *genericclioptionsapplier.ApplierFlags
		streams      genericclioptions.IOStreams
	}
	tests := []struct {
		name string
		args args
		want *Options
	}{
		{
			name: "success",
			args: args{
				applierFlags: applierFlags,
			},
			want: &Options{
				ApplierFlags: applierFlags,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newOptions(tt.args.applierFlags, tt.args.streams); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
