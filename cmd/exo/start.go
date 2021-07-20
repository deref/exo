package main

import (
	"github.com/deref/exo/exod/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start <ref>",
	Short: "Start a process",
	Long:  `Start a process.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		_, err := workspace.Start(ctx, &api.StartInput{
			Ref: args[0],
		})
		return err
	},
}
