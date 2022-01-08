package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringArrayVar(&lsFlags.Types, "type", nil, "filter by type")
}

var lsFlags struct {
	Types []string
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists components",
	Long:  `Lists components.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.DescribeComponents(ctx, &api.DescribeComponentsInput{
			Types: lsFlags.Types,
		})
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		for _, component := range output.Components {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", component.Name, component.ID, component.Type)
		}
		_ = w.Flush()
		return nil
	},
}
