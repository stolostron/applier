// Copyright Contributors to the Open Cluster Management project

//Package main: This program shows how to create resources based on yamls template located in
//the same directory or bindata or array of string.
package main

import (
	"flag"
	"os"

	"github.com/stolostron/applier/pkg/templateprocessor"
	"k8s.io/klog"
)

func usage() {
	klog.Info("Usage: render-yaml-in-dir\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	klog.InitFlags(nil)

	var showHelp = flag.Bool("h", false, "Show help message")

	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	err := renderYamlFile()
	if err != nil {
		klog.Fatal(err)
	}
}

func renderYamlFile() error {
	const directory = "../resources"
	//Create a reader on "../resources" directory
	klog.Infof("Creating the file reader %s", directory)
	yamlReader := templateprocessor.NewYamlFileReader(directory)
	//Other readers can be used
	//yamlReader := bindata.NewBindataReader()
	//yamlReader := templateprocessor.NewYamlStringReader(yamls,"---")

	//Create a templateProcessor with that reader
	klog.Infof("Creating TemplateProcessor...")
	tp, err := templateprocessor.NewTemplateProcessor(yamlReader, &templateprocessor.Options{})
	if err != nil {
		return err
	}
	//Defines the values
	values := struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}{
		ManagedClusterName:          "mycluster",
		ManagedClusterNamespace:     "mycluster",
		BootstrapServiceAccountName: "mybootstrapserviceaccount",
	}

	//Render the resources starting with path "yamlfilereader" in the reader
	//excluding "clusterrolebinding.yaml"
	//in a non-recursive way
	//and passing the values to replace
	//The output is NOT sorted, the order is defined by the templateNames parameter
	klog.Info("Render resources\n")
	out, err := tp.TemplateResources(
		[]string{"yamlfilereader/serviceaccount.yaml",
			"yamlfilereader/namespace.yaml"},
		values)
	if err != nil {
		return err
	}
	klog.Infof("Generated resources yamls\n%s", templateprocessor.ConvertArrayOfBytesToString(out))
	return nil
}
