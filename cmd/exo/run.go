package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&applyFlags.Format, "format", "", "see `exo help apply`")
}

var runFlags struct {
	Format string
}

var runCmd = &cobra.Command{
	Use:   "run [flags] [manifest-file]",
	Short: "Runs all processes and tails their logs",
	Long: `Runs all processes and tails their logs.

See 'exo help apply' for details on the manifest arguments.

If a workspace does not exist, one will be created in the current directory.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := newContext()
		ensureDaemon()
		cl := newClient()
		kernel := cl.Kernel()
		logger := logging.CurrentLogger(ctx)

		// Ensure workspace.
		workspace := mustFindWorkspace(ctx, cl)
		if workspace == nil {
			output, err := cl.Kernel().CreateWorkspace(ctx, &api.CreateWorkspaceInput{
				Root: cmdutil.MustGetwd(),
			})
			if err != nil {
				return fmt.Errorf("creating workspace: %w", err)
			}
			workspace = cl.GetWorkspace(output.ID)
		}

		// Advertise GUI URL.
		{
			output, err := workspace.Describe(ctx, &api.DescribeInput{})
			if err != nil {
				return fmt.Errorf("describing workspace: %w", err)
			}
			routes := newGUIRoutes()
			fmt.Println("GUI available at:", routes.WorkspaceURL(output.Description.ID))
		}

		// Apply manifest.
		if err := apply(ctx, kernel, workspace, args); err != nil {
			return fmt.Errorf("applying manifest: %w", err)
		}

		// Tail all logs until interrupt.
		(func() {
			ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()
			refs := []string{}
			stopOnError := true
			if err := tailLogs(ctx, workspace, refs, stopOnError); err != nil {
				logger.Infof("error tailing logs: %v", err)
			}
		})()

		// Stop workspace.
		fmt.Println("stopping workspace...")
		_, err := workspace.Stop(ctx, &api.StopInput{})
		if err != nil {
			return fmt.Errorf("stopping: %w", err)
		}
		fmt.Println("stopped")
		return nil
	},
}
