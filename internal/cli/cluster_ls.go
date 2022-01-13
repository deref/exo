package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

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
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		for _, cluster := range q.Clusters {
			labels := ""
			if cluster.Default {
				labels = "default"
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", cluster.ID, cluster.Name, labels)
		}
		_ = w.Flush()
		return nil
	},
}
