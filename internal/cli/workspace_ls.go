package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

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
		checkOrEnsureServer()

		cl, shutdown := dialGraphQL(ctx)
		defer shutdown()

		var q struct {
			Workspaces []struct {
				ID   string
				Root string
			} `graphql:"allWorkspaces"`
		}
		if err := cl.Query(ctx, &q, nil); err != nil {
			return fmt.Errorf("querying: %w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		fmt.Fprintln(w, "# ID\tROOT")
		for _, ws := range q.Workspaces {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", ws.ID, ws.Root)
		}
		_ = w.Flush()
		return nil
	},
}
