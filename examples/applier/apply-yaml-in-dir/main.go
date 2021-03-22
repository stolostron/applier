// Copyright Contributors to the Open Cluster Management project

//Package main: This program shows how to create resources based on yamls template located in
//the same directory or bindata or array of string.
package main

import (
	"flag"
	"os"

	// "github.com/open-cluster-management/applier/examples/applier/bindata"
	"github.com/open-cluster-management/applier/pkg/applier"
	"github.com/open-cluster-management/applier/pkg/templateprocessor"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func usage() {
	klog.Info("Usage: apply-yaml-in-dir -k kubeconfig\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	klog.InitFlags(nil)

	var kubeconfig = flag.String("k", "", "The path of the kubeconfig")
	var showHelp = flag.Bool("h", false, "Show help message")

	flag.Usage = usage
	flag.Parse()

	if *kubeconfig == "" {
		klog.Info("k is a mandatory argument")
		showUsageAndExit(0)
	}

	if *showHelp {
		showUsageAndExit(0)
	}

	err := applyYamlFile(*kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
}

func applyYamlFile(kubeconfig string) error {
	const directory = "../resources"
	//Create a reader on "../resources" directory
	klog.Infof("Creating the file reader %s", directory)
	yamlReader := templateprocessor.NewYamlFileReader(directory)
	//Other readers can be used
	//yamlReader := bindata.NewBindataReader()
	//yamlReader := templateprocessor.NewYamlStringReader(yamls,"---")

	//Create a client
	klog.Infof("Creating kubernetes client using kubeconfig located at %s", kubeconfig)
	apiconfig, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return err
	}
	config, err := clientcmd.NewDefaultClientConfig(
		*apiconfig,
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return err
	}
	client, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}
	//Create an Applier
	klog.Info("Creating applier")
	a, err := applier.NewApplier(yamlReader,
		&templateprocessor.Options{},
		client,
		nil,
		nil,
		applier.DefaultKubernetesMerger,
		nil)
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

	//Create the resources starting with path "yamlfilereader" in the reader
	//excluding "clusterrolebinding.yaml"
	//in a non-recursive way
	//and passing the values to replace
	klog.Info("Create or update resources")
	err = a.CreateOrUpdateInPath(
		"yamlfilereader",
		[]string{"yamlfilereader/clusterrolebinding.yaml"},
		false,
		values)
	if err != nil {
		return err
	}
	klog.Infof("Resource deployed")
	return nil
}
