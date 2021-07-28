package main

import (
	"github.com/deref/exo/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd)
}

var restartCmd = &cobra.Command{
	Use:   "restart [ref]",
	Short: "Restart a process",
	Long: `Restart a process. If it's not already running, will start it.

If a ref is not providered, restarts the entire workspace.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		switch len(args) {
		case 0:
			_, err := workspace.Restart(ctx, &api.RestartInput{})
			return err
		case 1:
			_, err := workspace.RestartComponent(ctx, &api.RestartComponentInput{
				Ref: args[0],
			})
			return err
		default:
			panic("unreachable")
		}
	},
}
