package cli

import (
	"github.com/deref/exo/internal/providers/docker/components/volume"
	"github.com/spf13/cobra"
)

func init() {
	newCmd.AddCommand(newVolumeCmd)
	// TODO: Options from docker to consider adding support for:
	// -d, --driver string   Specify volume driver name (default "local")
	// --label list      Set metadata for a volume
	// -o, --opt map         Set driver specific options (default map[])
}

var volumeSpec volume.Spec

var newVolumeCmd = &cobra.Command{
	Use:   "volume <name> [options]",
	Short: "Creates a new volume",
	Long: `Creates a new volume component.

Similar in spirit to:

docker volume create <name>
`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]
		return createComponent(ctx, name, "volume", volumeSpec)
	},
}
