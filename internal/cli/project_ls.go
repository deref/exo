package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	projectCmd.AddCommand(projectLSCmd)
}

var projectLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists projects",
	Long:  `Lists projects.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		gqlClient, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Projects []struct {
				ID          string
				DisplayName string
			} `graphql:"allProjects"`
		}
		if err := gqlClient.Query(ctx, &q, nil); err != nil {
			cmdutil.Fatalf("querying: %w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		for _, project := range q.Projects {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", project.ID, project.DisplayName)
		}
		_ = w.Flush()
		return nil
	},
}
