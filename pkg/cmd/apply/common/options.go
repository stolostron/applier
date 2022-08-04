// Copyright Contributors to the Open Cluster Management project
package common

import (
	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type ResourceType string

const CoreResources ResourceType = "core-resources"
const Deployments ResourceType = "deployments"
const CustomResources ResourceType = "custom-resources"

type Options struct {
	//ApplierFlags: The generic options from the applier cli-runtime.
	ApplierFlags *genericclioptionsapplier.ApplierFlags
	// Header specify a file that needs to be added at the beginning of each template
	Header string
	//A list of Paths
	Paths         []string
	ValuesPath    string
	Values        map[string]interface{}
	ResourcesType ResourceType
	//The file to output the resources will be sent to the file.
	OutputFile string
	SortOnKind bool
}

func NewOptions(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		ApplierFlags: applierFlags,
	}
}
