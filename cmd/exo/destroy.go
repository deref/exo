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
	Short: "Deletes all components in the project",
	Long:  `Deletes all components in the project.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDeamon()
		client := newClient()
		_, err := client.Destroy(ctx, &api.DestroyInput{})
		return err
	},
}
