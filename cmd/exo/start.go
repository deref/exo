package main

import (
	"github.com/deref/exo/kernel/api"
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
		ensureDeamon()
		client := newClient()
		_, err := client.Start(ctx, &api.StartInput{
			Ref: args[0],
		})
		return err
	},
}
