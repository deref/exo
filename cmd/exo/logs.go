package main

import (
	"fmt"
	"time"

	"github.com/deref/exo/exod/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logsCmd)
}

var logsCmd = &cobra.Command{
	Hidden: true,
	Use:    "logs [refs...]",
	Short:  "Tails process logs",
	Long: `Tails process logs.

If refs are provided, filters for the logs of those processes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		workspace := requireWorkspace(ctx, cl)
		cursor := ""
		for {
			output, err := workspace.GetEvents(ctx, &api.GetEventsInput{
				Logs:  args,
				After: cursor,
			})
			if err != nil {
				return err
			}
			for _, event := range output.Events {
				fmt.Println(event)
				fmt.Printf("%s %s %s\n", event.Log, event.Timestamp, event.Message)
			}
			cursor = output.Cursor
			if len(output.Events) < 10 { // TODO: OK heuristic?
				<-time.After(250 * time.Millisecond)
			}
		}
	},
}
