package cli

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/spf13/cobra"
)

func init() {
	manifestCmd.AddCommand(manifestFormatCmd)
}

var manifestFormatCmd = &cobra.Command{
	Use:   "format [<manifest>]",
	Short: "Formats a manifest file",
	Long: `Reformats a manifest file at the given path.

If the path is not provided, reads from stdin and writes to stdout.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		in := "/dev/stdin"
		out := "/dev/stdout"
		if len(args) > 0 {
			in = args[0]
			out = in
		}
		m, err := loadManifest(ctx, os.Stderr, in)
		if err != nil {
			return err
		}
		f, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return fmt.Errorf("opening out: %w", err)
		}
		defer f.Close()
		_, err = hclgen.WriteTo(f, hclgen.FileFromStructure(m.File))
		return err
	},
}
