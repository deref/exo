package cli

import (
	"fmt"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clusterCmd)

	clusterCmd.AddCommand(makeHelpSubcmd())
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Create, inspect, and modify clusters",
	Long: `Contains subcommands for operating on clusters.

If no subcommand is given, describes the cluster cluster.

The current cluster can be set explicitly with --cluster, otherwise it comes
from the current stack.  If there is no current stack, the default cluster is
used.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		type clusterFragment struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		}

		var cluster *clusterFragment
		if cmd.Flags().Lookup("cluster").Changed {
			var q struct {
				Cluster *clusterFragment `graphql:"clusterByRef(ref: $cluster)"`
			}
			if err := api.Query(ctx, svc, &q, map[string]interface{}{
				"cluster": rootPersistentFlags.Cluster,
			}); err != nil {
				return err
			}
			cluster = q.Cluster
			if cluster == nil {
				return fmt.Errorf("no such cluster: %q", rootPersistentFlags.Cluster)
			}
		} else {
			var q struct {
				Stack *struct {
					Cluster clusterFragment
				} `graphql:"stackByRef(ref: $stack)"`
				DefaultCluster *clusterFragment `graphql:"defaultCluster"`
			}
			if err := api.Query(ctx, svc, &q, map[string]interface{}{
				"stack": currentStackRef(),
			}); err != nil {
				return err
			}
			if q.Stack != nil {
				cluster = &q.Stack.Cluster
			} else if q.DefaultCluster != nil {
				cluster = q.DefaultCluster
			} else {
				return fmt.Errorf("no current cluster")
			}
		}

		cmdutil.PrintCueStruct(cluster)
		return nil
	},
}
