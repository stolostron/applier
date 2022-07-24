// Copyright Contributors to the Open Cluster Management project
package render

import (
	"fmt"

	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"github.com/stolostron/applier/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# render templates
%[1]s render --values values.yaml --path template_path1 tempalte_path2...
`

// NewCmd ...
func NewCmd(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(applierFlags, streams)

	cmd := &cobra.Command{
		Use:          "render",
		Short:        "render templates located in paths",
		Long:         "render templates located in paths with a values.yaml, the list of path can be a path to a file or a directory",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.ValuesPath, "values", "", "The files containing the values")
	cmd.Flags().StringVar(&o.Header, "header", "", "The files which will be added to each template")
	cmd.Flags().StringArrayVar(&o.Paths, "paths", []string{}, "The list of template paths")
	cmd.Flags().StringVar(&o.OutputFile, "output-file", "", "The generated resources will be copied in the specified file")
	return cmd
}
