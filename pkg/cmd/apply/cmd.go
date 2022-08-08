// Copyright Contributors to the Open Cluster Management project
package apply

import (
	"fmt"

	"github.com/stolostron/applier/pkg/cmd/apply/common"
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
%[1]s apply [core-resources|custom-resources|deployments] --values values.yaml --path template_path1 --path tempalte_path2...
`

// NewCmd ...
func NewCmd(applierFlags *genericclioptionsapplier.ApplierFlags, streams genericclioptions.IOStreams) *cobra.Command {
	o := common.NewOptions(applierFlags, streams)

	cmd := &cobra.Command{
		Use:          "apply",
		Short:        "apply templates located in paths",
		Long:         "apply templates located in paths with a values.yaml, the list of path can be a path to a file or a directory",
		Example:      fmt.Sprintf(example, helpers.GetExampleHeader()),
		SilenceUsage: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			helpers.DryRunMessage(o.ApplierFlags.DryRun)
		},
	}

	cmd.AddCommand(core.NewCmd(applierFlags, streams))
	cmd.AddCommand(customresources.NewCmd(applierFlags, streams))
	cmd.AddCommand(deployments.NewCmd(applierFlags, streams))
	return cmd
}
