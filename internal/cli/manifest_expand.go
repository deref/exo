package cli

import (
	"os"

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
		ctx := newContext()
		res, err := loadManifest(ctx, os.Stderr, args[0])
		if err != nil {
			return err
		}
		// XXX Print expanded version of manifest as HCL, not internal data structures.
		_, _ = pretty.Println(res)
		return nil
	},
}
