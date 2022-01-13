package cli

import (
	"fmt"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	clusterCmd.AddCommand(clusterLSCmd)
}

var clusterLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists clusters",
	Long:  `Lists clusters.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Clusters []struct {
				ID      string
				Name    string
				Default bool
			} `graphql:"allClusters"`
		}
		if err := cl.Query(ctx, &q, nil); err != nil {
			return fmt.Errorf("querying: %w", err)
		}
		w := cmdutil.NewTableWriter("NAME", "ID", "MISC")
		for _, cluster := range q.Clusters {
			misc := ""
			if cluster.Default {
				misc = "default"
			}
			w.WriteRow(cluster.Name, cluster.ID, misc)
		}
		w.Flush()
		return nil
	},
}
