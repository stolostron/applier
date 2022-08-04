// Copyright Contributors to the Open Cluster Management project
package render

import (
	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Options struct {
	// Header specify a file that needs to be added at the beginning of each template
	Header string
	//A list of Paths
	Paths      []string
	ValuesPath string
	Values     map[string]interface{}
	OutputFile string
	SortOnKind bool
	OutputDir  string
}

func NewOptions(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{}
}
