// Copyright Contributors to the Open Cluster Management project

package cmd

import (
	"flag"
	"os"

	"github.com/spf13/cobra"

	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"github.com/stolostron/applier/pkg/helpers"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/cmd/options"
	"k8s.io/kubectl/pkg/cmd/plugin"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/stolostron/applier/pkg/cmd/apply"
	"github.com/stolostron/applier/pkg/cmd/render"
	"github.com/stolostron/applier/pkg/cmd/version"
)

func NewCMCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "applier",
		Short: "apply templated resources",
		//This remove the auto-generated tag in the cobra doc
		DisableAutoGenTag: true,
	}

	flags := root.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.WrapConfigFn = setQPS
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

	matchVersionKubeConfigFlags.AddFlags(flags)

	klog.InitFlags(nil)
	flags.AddGoFlagSet(flag.CommandLine)

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	// From this point and forward we get warnings on flags that contain "_" separators
	root.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	applierFlags := genericclioptionsapplier.NewApplierFlags(f)
	applierFlags.AddFlags(flags)

	// root.AddCommand(cmdconfig.NewCmdConfig(f, clientcmd.NewDefaultPathOptions(), streams))
	root.AddCommand(options.NewCmdOptions(streams.Out))

	//enable plugin functionality: all `os.Args[0]-<binary>` in the $PATH will be available for plugin
	plugin.ValidPluginFilenamePrefixes = []string{os.Args[0]}
	root.AddCommand(plugin.NewCmdPlugin(streams))

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

func setQPS(r *rest.Config) *rest.Config {
	r.QPS = helpers.QPS
	r.Burst = helpers.Burst
	return r
}
