package cli

import (
	"fmt"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

func init() {
	workspaceCmd.AddCommand(workspaceLSCmd)
}

var workspaceLSCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists workspaces",
	Long:  `Lists workspaces.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var q struct {
			Workspaces []struct {
				ID   string
				Root string
			} `graphql:"allWorkspaces"`
		}
		if err := client.Query(ctx, &q, nil); err != nil {
			return fmt.Errorf("querying: %w", err)
		}
		w := cmdutil.NewTableWriter("ID", "ROOT")
		for _, ws := range q.Workspaces {
			w.WriteRow(ws.ID, ws.Root)
		}
		w.Flush()
		return nil
	},
}
