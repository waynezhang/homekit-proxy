package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/homekit-proxy/internal/constants"
	"github.com/waynezhang/homekit-proxy/internal/homekit"
)

func init() {
	var dbPath string
	var configFile string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Run: func(cmd *cobra.Command, args []string) {
			homekit.Serve(configFile, dbPath)
		},
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", "./db", "Database path")
	cmd.Flags().StringVarP(&configFile, "config", "c", constants.DefaultConfigFile, "Config file path")

	RootCmd.AddCommand(cmd)
}
