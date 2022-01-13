package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	resourceCmd.AddCommand(resourceNewCmd)
}

var resourceNewCmd = &cobra.Command{
	Use:   "new <type>",
	Short: "Create a component",
	Long:  "Create a new component.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("NOT YET IMPLEMENTED")
		//typ := args[0]

		//if term.IsInteractive() {
		//	spec, err := term.EditString("resource.*.txt", "")
		//}
	},
}
