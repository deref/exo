package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lists defined processes",
	Long:  `Describes defined processes and their statuses.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		output, err := workspace.DescribeProcesses(ctx, &api.DescribeProcessesInput{})
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		fmt.Fprintln(w, "# NAME\tID\tSTATE\tPROVIDER")
		for _, process := range output.Processes {
			state := "stopped"
			if process.Running {
				state = "running"
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", process.Name, process.ID, state, process.Provider)
		}
		_ = w.Flush()
		return nil
	},
}
