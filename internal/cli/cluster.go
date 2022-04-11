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
		cluster, err := lookupCluster(cmd)
		if err != nil {
			return err
		}
		cmdutil.PrintCueStruct(cluster)
		return nil
	},
}

type clusterFragment struct {
	ID          string
	Name        string
	Default     bool
	Environment environmentFragment
}

func lookupCluster(cmd *cobra.Command) (*clusterFragment, error) {
	ctx := cmd.Context()
	if cmd.Flags().Lookup("cluster").Changed {
		var q struct {
			Cluster *clusterFragment `graphql:"clusterByRef(ref: $cluster)"`
		}
		if err := api.Query(ctx, svc, &q, map[string]any{
			"cluster": rootPersistentFlags.Cluster,
		}); err != nil {
			return nil, err
		}
		if q.Cluster == nil {
			return nil, fmt.Errorf("no such cluster: %q", rootPersistentFlags.Cluster)
		}
		return q.Cluster, nil
	} else {
		var q struct {
			Stack *struct {
				Cluster clusterFragment
			} `graphql:"stackByRef(ref: $stack)"`
			DefaultCluster clusterFragment `graphql:"defaultCluster"`
		}
		if err := api.Query(ctx, svc, &q, map[string]any{
			"stack": currentStackRef(),
		}); err != nil {
			return nil, err
		}
		if q.Stack != nil {
			return &q.Stack.Cluster, nil
		} else {
			return &q.DefaultCluster, nil
		}
	}
}

func showCluster(cluster *clusterFragment) {
	env := make(map[string]any, len(cluster.Environment.Variables))
	for _, v := range cluster.Environment.Variables {
		env[v.Name] = v.Value
	}
	cmdutil.PrintCueStruct(map[string]any{
		"id":          cluster.ID,
		"name":        cluster.Name,
		"default":     cluster.Default,
		"environment": env,
	})
}
