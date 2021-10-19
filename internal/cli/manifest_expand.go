package cli

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestExpandCmd)
}

var manifestExpandCmd = &cobra.Command{
	Use:   "expand <manifest>",
	Short: "Loads a manifest file and prints it in expanded form.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := loadManifest(args[0])
		if err != nil {
			return fmt.Errorf("loading manifest: %w", err)
		}
		// XXX Print manifest as HCL, not internal data structures.
		_, _ = pretty.Print(res)
		return nil
	},
}
