package main

import (
	"github.com/deref/exo/exod/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop <ref>",
	Short: "Stop a process",
	Long:  `Stop a process.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		_, err := workspace.Stop(ctx, &api.StopInput{
			Ref: args[0],
		})
		return err
	},
}
