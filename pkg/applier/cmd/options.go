// Copyright Contributors to the Open Cluster Management project

package cmd

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	// "k8s.io/klog"
)

type Options struct {
	ConfigFlags *genericclioptions.ConfigFlags

	OutFile    string
	Directory  string
	ValuesPath string
	DryRun     bool
	Prefix     string
	Delete     bool
	Timeout    int
	Force      bool
	Silent     bool

	genericclioptions.IOStreams
}

func newOptions(streams genericclioptions.IOStreams) *Options {
	return &Options{
		ConfigFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}
