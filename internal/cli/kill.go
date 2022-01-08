package cli

import (
	"github.com/deref/exo/internal/core/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().StringVarP(&killFlags.Signal, "signal", "s", "kill", "")
}

var killFlags struct {
	Signal string
}

var killCmd = &cobra.Command{
	Use:   "kill [ref...]",
	Short: "Force stop processes",
	Long: `Sends a unix signal to process-like components.
	
If no refs are provided, signals all processes in the workspace.
	
The default signal is SIGKILL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)

		// We could use controlComponents here, but that does an implicit `watchJob`
		// which we don't want, since signal sending is so fast that we want to treat
		// it as synchronous.
		var err error
		if len(args) == 0 {
			_, err = workspace.Signal(ctx, &api.SignalInput{
				Signal: killFlags.Signal,
			})
		} else {
			_, err = workspace.SignalComponents(ctx, &api.SignalComponentsInput{
				Refs:   args,
				Signal: killFlags.Signal,
			})
		}
		return err
	},
}
