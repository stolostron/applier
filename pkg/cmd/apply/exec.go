// Copyright Red Hat
package apply

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/stolostron/applier/pkg/apply"
	"github.com/stolostron/applier/pkg/asset"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func (o *Options) Complete(cmd *cobra.Command, args []string) (err error) {
	// Convert yaml to map[string]interface
	b := []byte("")
	switch {
	case len(o.options.ValuesPath) == 0:
		// check if pipe
		fi, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		if fi.Mode()&os.ModeNamedPipe != 0 {
			b, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
		}
	case len(o.options.ValuesPath) != 0:
		b, err = ioutil.ReadFile(o.options.ValuesPath)
		if err != nil {
			return err
		}
	}
	o.options.Values = make(map[string]interface{})
	if err := yaml.Unmarshal(b, &o.options.Values); err != nil {
		return err
	}
	return nil
}

func (o *Options) Validate() error {
	reader := asset.NewDirectoriesReader(o.options.Header, o.options.Paths)

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
	kubeClient, err := o.options.ApplierFlags.KubectlFactory.KubernetesClientSet()
	if err != nil {
		return err
	}
	dynamicClient, err := o.options.ApplierFlags.KubectlFactory.DynamicClient()
	if err != nil {
		return err
	}

	restConfig, err := o.options.ApplierFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	apiExtensionsClient, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	applyBuilder := apply.NewApplierBuilder().
		WithClient(kubeClient, apiExtensionsClient, dynamicClient)
	reader := asset.NewDirectoriesReader(o.options.Header, o.options.Paths)
	files, err := reader.AssetNames([]string{o.options.Header})
	if err != nil {
		return err
	}
	if !o.options.SortOnKind {
		applyBuilder = applyBuilder.WithKindOrder(apply.NoCreateUpdateKindsOrder)
	}
	applier := applyBuilder.Build()
	output, err := applier.Apply(reader, o.options.Values, o.options.ApplierFlags.DryRun, o.options.Header, files...)
	if err != nil {
		return err
	}
	return apply.WriteOutput(o.options.OutputFile, output)
}
