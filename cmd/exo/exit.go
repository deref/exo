package main

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exitCmd)
}

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Stop the exo daemon",
	Long:  `Stop the exo daemon process.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		_, err := cl.Kernel().Exit(ctx, &api.ExitInput{})
		if err != nil {
			cmdutil.Fatalf("exiting: %w", err)
		}
		fmt.Println("Ok")
		return nil
	},
}
