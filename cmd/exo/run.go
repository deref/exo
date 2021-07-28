package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/deref/exo/core/api"
	"github.com/deref/exo/util/cmdutil"
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
			fmt.Println("GUI available at:", guiWorkspaceURL(output.Description.ID))
		}

		// Apply manifest.
		if err := apply(ctx, workspace, args); err != nil {
			return fmt.Errorf("applying manifest: %w", err)
		}

		// Tail all logs until interrupt.
		(func() {
			ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
			defer stop()
			var logRefs []string
			if err := tailLogs(ctx, workspace, logRefs); err != nil {
				log.Printf("error tailing logs: %v", err)
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
