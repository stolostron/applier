// Copyright Contributors to the Open Cluster Management project
package version

import (
	"fmt"

	"github.com/spf13/cobra"
	appliercli "github.com/stolostron/applier"
)

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	return nil
}

func (o *Options) validate() error {
	return nil
}
func (o *Options) run() (err error) {
	fmt.Printf("client\t\t\tversion\t:%s\n", appliercli.GetVersion())
	return nil
}
