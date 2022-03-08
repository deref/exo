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

		type versionInfo struct {
			Version string `json:"version" graphql:"installed"`
			Build   string `json:"build"`
			Managed bool   `json:"managed"`
		}
		client := versionInfo{
			Version: about.Version,
			Build:   about.GetBuild(),
			Managed: about.IsManaged,
		}

		res := map[string]interface{}{}

		if isPeerMode() {
			res["peer"] = client
		} else {
			res["client"] = client
			var q struct {
				System struct {
					Version versionInfo
				}
			}
			if err := api.Query(ctx, svc, &q, nil); err != nil {
				cmdutil.Warnf("getting server version: %v", err)
			} else {
				res["server"] = q.System.Version
			}
		}

		cmdutil.PrintCueStruct(res)
	},
}
