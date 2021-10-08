package cli

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/manifest"
	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestValidateCmd)
}

var manifestValidateCmd = &cobra.Command{
	Use:   "validate <manifest>",
	Short: "Loads and validates a manifest file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Support other formats here too.
		loader := &manifest.Loader{}
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("opening manifest: %w", err)
		}
		defer f.Close()
		res := loader.Load(f)
		for _, warning := range res.Warnings {
			fmt.Println(warning)
		}
		return res.Err
	},
}
