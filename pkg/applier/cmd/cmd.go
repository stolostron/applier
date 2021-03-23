// Copyright Contributors to the Open Cluster Management project

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	// "k8s.io/klog"
)

var example = `
# Apply templates
%[1]s applier -d <templatepath> --values values.yaml
`

//NewCmd generates a cobra.Command
func NewCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := newOptions(streams)

	cmd := &cobra.Command{
		Use:          "applier",
		Short:        "apply templates",
		Example:      fmt.Sprintf(example, getExampleHeader()),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.complete(c, args); err != nil {
				return err
			}
			if err := o.validate(); err != nil {
				return err
			}
			if err := o.run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.OutFile, "output", "o", "",
		"Output file. If set nothing will be applied but a file will be generate "+
			"which you can apply later with 'kubectl <create|apply|delete> -f")
	cmd.Flags().StringVarP(&o.Directory, "directory", "d", "", "The directory or file containing the template(s).\n"+
		"If a `_helpers.tpl` file exists in the same directory of the file, the `_helpers.tpl` will be included.")
	cmd.Flags().StringVar(&o.ValuesPath, "values", "", "The file containing the values")
	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "if set only the rendered yaml will be shown, default false")
	cmd.Flags().StringVarP(&o.Prefix, "prefix", "p", "", "The prefix to add to each value names, for example 'Values'")
	cmd.Flags().BoolVar(&o.Delete, "delete", false,
		"if set only the resource defined in the yamls will be deleted, default false")
	cmd.Flags().IntVar(&o.Timeout, "timout", 5, "Timeout in second to apply one resource, default 5 sec")
	cmd.Flags().BoolVarP(&o.Force, "force", "f", false, "If set, the finalizers will be removed before delete")
	cmd.Flags().BoolVar(&o.Silent, "silent", false, "If set the applier will run silently")

	o.ConfigFlags.AddFlags(cmd.Flags())

	return cmd
}
