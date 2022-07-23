// Copyright Contributors to the Open Cluster Management project
package render

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
	if len(o.OutputFile) == 0 {
		o.OutputFile = os.Stdout.Name()
	}
	return nil
}

func (o *Options) Validate() error {
	reader := asset.NewDirectoriesReader(o.Header, o.Paths)

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
	applyBuilder := apply.NewApplierBuilder()
	applier := applyBuilder.Build()
	reader := asset.NewDirectoriesReader(o.Header, o.Paths)
	files, err := reader.AssetNames([]string{o.Header})
	if err != nil {
		return err
	}
	output, err := applier.MustTemplateAssets(reader, o.Values, o.Header, files...)
	if err != nil {
		return err
	}
	return apply.WriteOutput(o.OutputFile, output)
}
