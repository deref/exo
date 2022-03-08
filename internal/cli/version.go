package cli

import (
	"github.com/deref/exo/internal/about"
	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of exo",
	Long:  "Print the version of exo.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		res := map[string]interface{}{
			"client": about.GetVersionInfo(),
		}

		if isClientMode() {
			var q struct {
				System struct {
					Version string `json:"version"`
					Build   string `json:"build"`
				}
			}
			if err := api.Query(ctx, svc, &q, nil); err != nil {
				cmdutil.Warnf("getting server version: %v", err)
			} else {
				res["server"] = q.System
			}
		}

		cmdutil.PrintCueStruct(res)
	},
}
