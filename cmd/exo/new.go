package main

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new <type> <name> [args]",
	Short: "Creates a new component.",
	Long: `Creates a new component of a given type with a given name.  Each
component type may define its own syntax for flags and positional arguments.

To learn about specific types - for example Processes - consult each type's
help page:

exo help new process
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}
