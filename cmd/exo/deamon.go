package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deamonCmd)
}

var deamonCmd = &cobra.Command{
	Hidden: true,
	Use:    "deamon",
	Short:  "Start the exo deamon",
	Long: `Start the exo deamon and then do nothing else.

Since most commands implicitly start the exo deamon, users generally do not
have to invoke this themselves.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ensureDeamon()
		return nil
	},
}

func ensureDeamon() {
}
