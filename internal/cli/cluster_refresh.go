package cli

import (
	"github.com/deref/exo/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	clusterCmd.AddCommand(clusterRefreshCmd)
}

var clusterRefreshCmd = &cobra.Command{
	Hidden: true,
	Use:    "refresh",
	Short:  "Refresh cluster model",
	Long: `Forces a refreshes of a cluster's model.

Targets a cluster as per the root "cluster" command.

It is not normally necessary to call this command, since cluster models are
automatically refreshed at some periodic frequency.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cluster, err := lookupCluster(cmd)
		if err != nil {
			return err
		}
		var m struct {
			Cluster *clusterFragment `graphql:"refreshCluster(ref: $cluster)"`
		}
		if err := api.Mutate(ctx, svc, &m, map[string]any{
			"cluster": cluster.ID,
		}); err != nil {
			return err
		}
		showCluster(m.Cluster)
		return nil
	},
}
