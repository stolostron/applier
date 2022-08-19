// Copyright Red Hat
package apply

import (
	"fmt"

	"github.com/stolostron/applier/pkg/cmd/apply/core"
	"github.com/stolostron/applier/pkg/cmd/apply/customresources"
	"github.com/stolostron/applier/pkg/cmd/apply/deployments"
	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"github.com/stolostron/applier/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Apply templates
%[1]s apply --values values.yaml --path template_path1 --path tempalte_path2...
`

// NewCmd ...
func NewCmd(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(applierFlags, streams)

	cmd := &cobra.Command{
		Use:          "apply",
		Short:        "apply templates located in paths",
		Long:         "apply templates located in paths with a values.yaml, the list of path can be a path to a file or a directory",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			helpers.DryRunMessage(o.options.ApplierFlags.DryRun)
		},
		RunE: func(c *cobra.Command, args []string) error {
			return o.runE(c, args)
		},
	}

	cmd.Flags().BoolVar(&o.options.ApplierFlags.DryRun, "dry-run", false, "If set the resources will not be applied")
	cmd.Flags().StringVar(&o.options.Header, "header", "", "The files which will be added to each template")
	cmd.Flags().StringVar(&o.options.OutputFile, "output-file", "", "The generated resources will be copied in the specified file")
	cmd.Flags().StringVar(&o.options.ValuesPath, "values", "", "The files containing the values")
	cmd.Flags().StringArrayVar(&o.options.Paths, "path", []string{}, "The list of template paths")
	cmd.Flags().StringArrayVar(&o.options.Exclude, "exclude", []string{}, "The list of paths to exclude")
	cmd.Flags().BoolVar(&o.options.SortOnKind, "sort-on-kind", true, "If set the files will be sorted by their kind (default true)")

	cmd.AddCommand(core.NewCmd(applierFlags, streams))
	cmd.AddCommand(customresources.NewCmd(applierFlags, streams))
	cmd.AddCommand(deployments.NewCmd(applierFlags, streams))
	return cmd
}

func (o *Options) runE(c *cobra.Command, args []string) error {
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
}
