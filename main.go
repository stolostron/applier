// Copyright Red Hat

package main

import (
	"os"

	"k8s.io/klog/v2"

	"github.com/stolostron/applier/pkg/cmd"
)

func main() {
	root := cmd.NewApplierCommand()
	err := root.Execute()
	if err != nil {
		klog.V(1).ErrorS(err, "Error:")
	}
	klog.Flush()
	if err != nil {
		os.Exit(1)
	}
}
