// Copyright Contributors to the Open Cluster Management project

package main

import (
	"os"

	"github.com/spf13/pflag"

	appliercmd "github.com/open-cluster-management/applier/pkg/applier/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	flags := pflag.NewFlagSet("applier", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := appliercmd.NewCmd(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
