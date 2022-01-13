package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

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

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Projects []struct {
				ID          string
				DisplayName string
			} `graphql:"allProjects"`
		}
		if err := cl.Query(ctx, &q, nil); err != nil {
			return fmt.Errorf("querying: %w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		fmt.Fprintln(w, "# ID\tDISPLAY NAME")
		for _, project := range q.Projects {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", project.ID, project.DisplayName)
		}
		_ = w.Flush()
		return nil
	},
}
