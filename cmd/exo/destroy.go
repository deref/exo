package main

import (
	"github.com/deref/exo/kernel/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Deletes the current workspace",
	Long:  `Deletes all components in the current workspace, then the workspace itself.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDeamon()
		client := newClient()
		_, err := client.Destroy(ctx, &api.DestroyInput{})
		return err
	},
}
