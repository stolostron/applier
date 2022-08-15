// Copyright Red Hat

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/stolostron/applier/pkg/helpers"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/stolostron/applier/pkg/cmd/apply"
	"github.com/stolostron/applier/pkg/cmd/render"
	"github.com/stolostron/applier/pkg/cmd/version"
)

func NewApplierCommand() *cobra.Command {
	root, applierFlags, streams := helpers.NewRootCmd()
	groups := templates.CommandGroups{
		{
			Message: "General commands:",
			Commands: []*cobra.Command{
				version.NewCmd(applierFlags, streams),
				apply.NewCmd(applierFlags, streams),
				render.NewCmd(applierFlags, streams),
			},
		},
	}
	groups.Add(root)
	return root
}
