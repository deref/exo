package cli

import (
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
					ID          string `json:"id"`
					DisplayName string `json:"displayName"`
				}
			} `graphql:"workspaceByRef(ref: $currentWorkspace)"`
		}
		mustQueryWorkspace(ctx, &q, nil)
		cmdutil.PrintCueStruct(q.Workspace.Project)
		return nil
	},
}

func currentProjectRef() string {
	// Allow override with a persistent flag or other non working-directory state.
	return cmdutil.MustGetwd()
}
