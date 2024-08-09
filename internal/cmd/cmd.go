package cmd

import (
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
	"github.com/spf13/cobra"
)

var RootCmd = func() *cobra.Command {
	updateLogger(slog.LevelWarn)

	var verbose bool
	var rootCmd = &cobra.Command{
		Use:   "homekit-proxy",
		Short: "Homekit Proxy",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				updateLogger(slog.LevelInfo)
			}
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	return rootCmd
}()

func updateLogger(level slog.Leveler) {
	slog.SetDefault(slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{Level: level})))
}
