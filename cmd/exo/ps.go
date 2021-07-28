package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/deref/exo/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Lists defined processes",
	Long:  `Describes defined processes and their statuses.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		output, err := workspace.DescribeProcesses(ctx, &api.DescribeProcessesInput{})
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0)
		for _, process := range output.Processes {
			state := "stopped"
			if process.Running {
				state = "running"
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", process.Name, process.ID, state)
		}
		_ = w.Flush()
		return nil
	},
}
