package cli

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Shows Exo configuration",
	Long:  `Shows Exo configuration.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// XXX Cue, not TOML.
		enc := toml.NewEncoder(os.Stdout)
		return enc.Encode(cfg)
	},
}
