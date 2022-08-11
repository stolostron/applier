// Copyright Red Hat
package version

import (
	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	//ApplierFlags: The generic optiosn from the cm cli-runtime.
	ApplierFlags *genericclioptionsapplier.ApplierFlags
}

func newOptions(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		ApplierFlags: applierFlags,
	}
}
