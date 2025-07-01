package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/homekit-proxy/internal/homekit"
)

func init() {
	var dbPath string
	var configDir string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Run: func(cmd *cobra.Command, args []string) {
			homekit.Serve(configDir, dbPath)
		},
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", "./db", "Database path")
	cmd.Flags().StringVarP(&configDir, "config", "c", ".", "Config directory path")

	RootCmd.AddCommand(cmd)
}
