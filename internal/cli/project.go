package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/internal/util/cmdutil"
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

If no subcommand is given, describes the current project.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var q struct {
			Workspace *struct {
				ID      string
				Project struct {
					ID          string
					DisplayName string
				}
			} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
		}
		mustQueryWorkspace(ctx, &q, nil)
		project := q.Workspace.Project
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "id:\t%s\n", project.ID)
		_, _ = fmt.Fprintf(w, "display-name:\t%s\n", project.DisplayName)
		_ = w.Flush()
		return nil
	},
}

func currentProjectRef() string {
	// Allow override with a persistent flag or other non working-directory state.
	return cmdutil.MustGetwd()
}
