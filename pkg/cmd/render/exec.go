// Copyright Red Hat
package render

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
		if fi.Mode()&os.ModeNamedPipe != 0 {
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

	if len(o.OutputFile) == 0 {
		o.OutputFile = os.Stdout.Name()
	}
	return nil
}

func (o *Options) Validate() error {
	reader := asset.NewDirectoriesReader(o.Header, o.Paths)

	assetNames, err := reader.AssetNames(o.Paths, []string{o.Header})
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
	if !o.SortOnKind {
		applyBuilder = applyBuilder.WithKindOrder(apply.NoCreateUpdateKindsOrder)
	}
	applier := applyBuilder.Build()

	reader := asset.NewDirectoriesReader(o.Header, o.Paths)
	// Get files names
	files, err := reader.AssetNames(o.Paths, []string{o.Header})
	if err != nil {
		return err
	}

	if len(o.OutputDir) == 0 {
		output, err := applier.MustTemplateAssets(reader, o.Values, o.Header, files...)
		if err != nil {
			return err
		}
		return apply.WriteOutput(o.OutputFile, output)
	} else {
		for _, name := range files {
			if name == o.Header {
				continue
			}
			output, err := applier.MustTemplateAsset(reader, o.Values, o.Header, name)
			if err != nil {
				return err
			}
			newFileName := filepath.Join(o.OutputDir, name)
			err = os.MkdirAll(filepath.Dir(newFileName), 0700)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(newFileName, output, 0600)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
