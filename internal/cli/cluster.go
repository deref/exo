package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/internal/api"
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
			ID      string
			Name    string
			Default bool
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

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", cluster.ID)
		_, _ = fmt.Fprintf(w, "name:\t%s\n", cluster.Name)
		_, _ = fmt.Fprintf(w, "default:\t%v\n", cluster.Default)
		_ = w.Flush()
		return nil
	},
}
