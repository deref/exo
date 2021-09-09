package main

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker/components/volume"
	"github.com/deref/exo/internal/util/yamlutil"
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
	Long: `Creates a new volume.

Similar in spirit to:

docker volume create <name>
`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)

		name := args[0]

		output, err := workspace.CreateComponent(ctx, &api.CreateComponentInput{
			Name: name,
			Type: "volume",
			Spec: yamlutil.MustMarshalString(volumeSpec),
		})
		if err != nil {
			return err
		}
		fmt.Println(output.ID)
		return nil
	},
}
