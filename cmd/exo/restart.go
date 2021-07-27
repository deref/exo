package main

import (
	"github.com/deref/exo/exod/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd)
}

var restartCmd = &cobra.Command{
	Use:   "restart <ref>",
	Short: "Restart a process",
	Long:  `Restart a process. If it's not already running, will start it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		_, err := workspace.RestartComponent(ctx, &api.RestartComponentInput{
			Ref: args[0],
		})
		return err
	},
}
