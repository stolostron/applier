// Copyright Red Hat
package deployments

import (
	"fmt"

	"github.com/stolostron/applier/pkg/cmd/apply/common"
	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"github.com/stolostron/applier/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Apply deployments templates
%[1]s apply deployments --values values.yaml --path template_path1 --path tempalte_path2...
`

// NewCmd ...
func NewCmd(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := common.NewOptions(applierFlags, streams)

	cmd := &cobra.Command{
		Use:          "deployments",
		Short:        "apply deployments templates located in paths",
		Long:         "apply deployments templates located in paths with a values.yaml, the list of path can be a path to a file or a directory",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			o.ResourcesType = common.Deployments
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
	cmd.Flags().StringArrayVar(&o.Paths, "path", []string{}, "The list of template paths")
	cmd.Flags().StringArrayVar(&o.Excluded, "excluded", []string{}, "The list of paths to exclude")
	cmd.Flags().BoolVar(&o.ApplierFlags.DryRun, "dry-run", false, "If set the generated resources will be displayed but not applied")
	cmd.Flags().IntVar(&o.ApplierFlags.Timeout, "timeout", 300, "extend timeout from 300 secounds ")
	cmd.Flags().StringVar(&o.OutputFile, "output-file", "", "The generated resources will be copied in the specified file")
	return cmd
}
