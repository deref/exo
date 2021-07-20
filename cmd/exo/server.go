package main

import (
	"github.com/deref/exo/components/process"
	"github.com/deref/exo/exod"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Hidden: true,
	Use:    "server",
	Short:  "Runs the exo server",
	Long: `Runs the exo server until interrupted.

Prefer the daemonize command for normal operation.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		exod.Main(exod.Config{
			Fifofum: fifofumBundleConfig,
		})
	},
}

var fifofumBundleConfig = process.FifofumConfig{
	Path: "exo",
	Args: []string{"fifofum"},
}
