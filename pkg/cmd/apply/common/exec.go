// Copyright Contributors to the Open Cluster Management project
package common

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/stolostron/applier/pkg/apply"
	"github.com/stolostron/applier/pkg/asset"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func (o *Options) Complete(cmd *cobra.Command, args []string) (err error) {
	// Convert yaml to map[string]interface
	if len(o.ValuesPath) != 0 {
		b, err := ioutil.ReadFile(o.ValuesPath)
		if err != nil {
			return err
		}
		o.Values = make(map[string]interface{})
		if err := yaml.Unmarshal(b, &o.Values); err != nil {
			return err
		}
	}
	return nil
}

func (o *Options) Validate() error {
	reader := asset.NewDirectoriesReader(o.Header, o.Path)

	assetNames, err := reader.AssetNames(nil)
	if err != nil {
		return err
	}
	if len(assetNames) == 0 {
		return fmt.Errorf("no files selected")
	}
	return nil
}

func (o *Options) Run() error {
	kubeClient, err := o.ApplierFlags.KubectlFactory.KubernetesClientSet()
	if err != nil {
		return err
	}
	dynamicClient, err := o.ApplierFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}

	restConfig, err := o.ApplierFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	apiExtensionsClient, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	applyBuilder := apply.NewApplierBuilder().
		WithClient(kubeClient, apiExtensionsClient, dynamicClient)
	reader := asset.NewDirectoriesReader(o.Header, o.Path)
	files, err := reader.AssetNames([]string{o.Header})
	if err != nil {
		return err
	}
	output := make([]string, 0)
	switch o.ResourcesType {
	case CoreResources:
		if o.SortOnKind {
			applyBuilder = applyBuilder.WithKindOrder(apply.DefaultCreateUpdateKindsOrder)
		}
		applier := applyBuilder.Build()
		output, err = applier.ApplyDirectly(reader, o.Values, o.ApplierFlags.DryRun, o.Header, files...)
	case Deployments:
		applier := applyBuilder.Build()
		output, err = applier.ApplyDeployments(reader, o.Values, o.ApplierFlags.DryRun, o.Header, files...)
	case CustomResources:
		applier := applyBuilder.Build()
		output, err = applier.ApplyCustomResources(reader, o.Values, o.ApplierFlags.DryRun, o.Header, files...)
	}
	if err != nil {
		return err
	}
	return apply.WriteOutput(o.OutputFile, output)
}
