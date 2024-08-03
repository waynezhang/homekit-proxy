package cmd

import (
	"github.com/spf13/cobra"
	"github.com/waynezhang/homekit-proxy/internal/config"
	"github.com/waynezhang/homekit-proxy/internal/constants"
	"github.com/waynezhang/homekit-proxy/internal/homekit"
)

var ServeCmd = func() *cobra.Command {
	var dbPath string
	var configFile string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Run: func(cmd *cobra.Command, args []string) {
			serve(configFile, dbPath)
		},
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", "./db", "Database path")
	cmd.Flags().StringVarP(&configFile, "config", "c", constants.DefaultConfigFile, "Config file path")

	return cmd
}()

func serve(file string, dbPath string) {
	config := config.Parse(file)
	hm := homekit.New(&config)
	hm.Start(dbPath)
}
