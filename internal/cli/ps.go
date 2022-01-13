package cli

import (
	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/cmdutil"
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
		w := cmdutil.NewTableWriter("NAME", "ID", "STATE", "PROVIDER")
		for _, process := range output.Processes {
			state := "stopped"
			if process.Running {
				state = "running"
			}
			w.WriteRow(process.Name, process.ID, state, process.Provider)
		}
		w.Flush()
		return nil
	},
}
