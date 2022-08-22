// Copyright Red Hat
package common

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/stolostron/applier/pkg/apply"
	"github.com/stolostron/applier/pkg/asset"
)

func (o *Options) Complete(cmd *cobra.Command, args []string) (err error) {
	// Convert yaml to map[string]interface
	b := []byte("")
	switch {
	case len(o.ValuesPath) == 0:
		// check if pipe
		fi, err := os.Stdin.Stat()
		if err != nil {
			return err
		}
		if fi.Mode()&os.ModeCharDevice == 0 {
			b, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
		}
	case len(o.ValuesPath) != 0:
		b, err = ioutil.ReadFile(o.ValuesPath)
		if err != nil {
			return err
		}
	}
	o.Values = make(map[string]interface{})
	if err := yaml.Unmarshal(b, &o.Values); err != nil {
		return err
	}
	return nil
}

func (o *Options) Validate() error {
	reader, err := asset.NewDirectoriesReader(o.Header, o.Paths)
	if err != nil {
		return err
	}

	assetNames, err := reader.AssetNames(o.Paths, nil, o.Header)
	if err != nil {
		return err
	}
	if len(assetNames) == 0 {
		return fmt.Errorf("no files selected")
	}
	return nil
}

func (o *Options) Run() error {
	restConfig, err := o.ApplierFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	applyBuilder := apply.NewApplierBuilder().WithRestConfig(restConfig)
	reader, err := asset.NewDirectoriesReader(o.Header, o.Paths)
	if err != nil {
		return err
	}
	o.Exclude = append(o.Exclude, o.Header)
	files, err := reader.AssetNames(o.Paths, o.Exclude, o.Header)
	if err != nil {
		return err
	}
	output := make([]string, 0)
	switch o.ResourcesType {
	case CoreResources:
		if !o.SortOnKind {
			applyBuilder = applyBuilder.WithKindOrder(apply.NoCreateUpdateKindsOrder)
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
