package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = ""
var Revision = ""

func init() {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("HomeKit Proxy v%s+%s\n", Version, Revision)
		},
	}
	RootCmd.AddCommand(cmd)
}
