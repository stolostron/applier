// Copyright Contributors to the Open Cluster Management project
package version

import (
	"fmt"

	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"github.com/stolostron/applier/pkg/helpers"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var example = `
# Version
%[1]s version
`

// NewCmd provides a cobra command wrapping NewCmdImportCluster
func NewCmd(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) (cmd *cobra.Command) {
	o := newOptions(applierFlags, streams)
	cmd = &cobra.Command{
		Use:          "version",
		Short:        "get the versions of the different components",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
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

	return cmd
}
