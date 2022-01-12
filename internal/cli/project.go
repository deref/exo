package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(projectCmd)

	projectCmd.AddCommand(makeHelpSubcmd())
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Create, inspect, and modify projects",
	Long: `Contains subcommands for operating on projects.

If no subcommand is given, describes the project of the current workspace.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Workspace *struct {
				ID      string
				Project *struct {
					ID          string
					DisplayName string
				}
			} `graphql:"workspaceByRef(ref: $workspace)"`
		}
		mustQueryWorkspace(ctx, cl, &q, nil)
		project := q.Workspace.Project
		if project == nil {
			return fmt.Errorf("no project for workspace %q", q.Workspace.ID)
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", project.ID)
		_, _ = fmt.Fprintf(w, "display-name:\t%s\n", project.DisplayName)
		_ = w.Flush()
		return nil
	},
}