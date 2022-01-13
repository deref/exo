package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shurcooL/graphql"
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

If no subcommand is given, describes the cluster of the current stack.  If
there is no current stack, describes the default cluster.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		type clusterFragment struct {
			ID      string
			Name    string
			Default bool
		}
		var q struct {
			DefaultCluster *clusterFragment `graphql:"defaultCluster"`
			Stack          *struct {
				Cluster clusterFragment
			} `graphql:"stackByRef(ref: $stack)"`
		}
		if err := cl.Query(ctx, &q, map[string]interface{}{
			"stack": graphql.String(currentStackRef()),
		}); err != nil {
			return err
		}

		var cluster *clusterFragment
		if q.Stack != nil {
			cluster = &q.Stack.Cluster
		} else if q.DefaultCluster != nil {
			cluster = q.DefaultCluster
		} else {
			return fmt.Errorf("no current cluster")
		}

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", cluster.ID)
		_, _ = fmt.Fprintf(w, "name:\t%s\n", cluster.Name)
		_, _ = fmt.Fprintf(w, "default:\t%v\n", cluster.Default)
		_ = w.Flush()
		return nil
	},
}
