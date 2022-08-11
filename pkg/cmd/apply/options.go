// Copyright Red Hat
package apply

import (
	"github.com/stolostron/applier/pkg/cmd/apply/common"
	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type ResourceType string

type Options struct {
	options common.Options
}

func NewOptions(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		options: common.Options{
			ApplierFlags: applierFlags,
		},
	}
}
