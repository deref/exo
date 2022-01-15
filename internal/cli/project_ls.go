package cli

import (
	"fmt"

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

		var q struct {
			Projects []struct {
				ID          string
				DisplayName string
			} `graphql:"allProjects"`
		}
		if err := client.Query(ctx, &q, nil); err != nil {
			return fmt.Errorf("querying: %w", err)
		}
		w := cmdutil.NewTableWriter("ID", "DISPLAY NAME")
		for _, project := range q.Projects {
			w.WriteRow(project.ID, project.DisplayName)
		}
		w.Flush()
		return nil
	},
}
