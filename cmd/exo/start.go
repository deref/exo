package main

import (
	"github.com/deref/exo/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start [ref]",
	Short: "Start a process",
	Long: `Start a process.

If a ref is not provided, starts the entire workspace.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		switch len(args) {
		case 0:
			_, err := workspace.Start(ctx, &api.StartInput{})
			return err
		case 1:
			_, err := workspace.StartComponent(ctx, &api.StartComponentInput{
				Ref: args[0],
			})
			return err
		default:
			panic("unreachable")
		}
	},
}
