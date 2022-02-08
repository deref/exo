package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Nerdmaster/terminal"
	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/core/api"
	joshclient "github.com/deref/exo/internal/core/client"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/mathutil"
	"github.com/deref/exo/internal/util/term"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolVarP(&logFlags.System, "system", "", false, "if specified, filter includes workspace system events")
	logsCmd.Flags().BoolVarP(&logFlags.NoFollow, "no-follow", "", false, "stops tailing when caught up to the end of the log")
}

var logFlags struct {
	System   bool
	NoFollow bool
}

var logsCmd = &cobra.Command{
	Hidden: true,
	Use:    "logs [flags] [refs...]",
	Short:  "Tails process logs",
	Long: `Tails and follows process logs.

If refs are provided, filters for the logs of those processes.

When filtering, system events are omitted unless --system is given.`,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		stopOnError := false
		return tailLogs(ctx, workspace, args, stopOnError)
	},
}

func tailLogs(ctx context.Context, workspace *joshclient.Workspace, streamRefs []string, stopOnError bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var eg errgroup.Group
	interactive := !logFlags.NoFollow
	if interactive {
		eg.Go(func() error {
			return runTailLogsReader(ctx, cancel)
		})
	}
	eg.Go(func() error {
		return runTailLogsWriter(ctx, workspace, streamRefs, stopOnError)
	})
	return eg.Wait()
}

func runTailLogsReader(ctx context.Context, cancel func()) error {
	if !term.IsInteractive() {
		return nil
	}

	raw := &term.RawMode{}
	raw.Enter()
	defer func() {
		if err := raw.Exit(); err != nil {
			cmdutil.Fatalf("restoring terminal state: %w", err)
		}
	}()

	r := terminal.NewKeyReader(os.Stdin)

	for {
		press, err := r.ReadKeypress()
		if err != nil {
			return fmt.Errorf("reading terminal keypress: %w", err)
		}

		switch press.Key {
		// Clear screen.
		case terminal.KeyCtrlL:
			fmt.Print("\033[2J\033[1;1H")

		// Quit.
		case terminal.KeyCtrlC, 'q':
			cancel()
			return nil

		// Append blank lines.
		// Simulates the common readline behavior to allow people to insert
		// a visual separator in their logs output.
		case '\r':
			fmt.Print("\r\n")

		// Suspend.
		case terminal.KeyCtrlZ:
			if err := raw.Suspend(); err != nil {
				cmdutil.Fatalf("suspending: %w", err)
			}
		}
	}
}

func runTailLogsWriter(ctx context.Context, workspace *joshclient.Workspace, streamRefs []string, stopOnError bool) error {
	workspaceID := workspace.ID()

	showName := len(streamRefs) != 1

	resolved, err := workspace.Resolve(ctx, &api.ResolveInput{
		Refs: streamRefs,
	})
	if err != nil {
		return fmt.Errorf("resolving refs: %w", err)
	}
	var streamNames []string
	if logFlags.System || len(streamRefs) > 0 {
		streamNames = make([]string, 0, 1+len(streamRefs))
		if logFlags.System {
			streamNames = append(streamNames, workspaceID)
		}
		for _, logID := range resolved.IDs {
			if logID != nil {
				streamNames = append(streamNames, *logID)
			}
		}
	}

	// TODO: Listen to change events to handle renames, new processes, etc.
	descriptions, err := workspace.DescribeProcesses(ctx, &api.DescribeProcessesInput{})
	if err != nil {
		return fmt.Errorf("describing processes: %w", err)
	}

	el := &EventLogger{
		W: os.Stdout,
	}

	streamToLabel := make(map[string]string, len(descriptions.Processes))
	streamToLabel[workspaceID] = "EXO"
	for _, process := range descriptions.Processes {
		streamToLabel[process.ID] = process.Name
		el.LabelWidth = mathutil.IntMax(el.LabelWidth, len(process.Name))
	}

	limit := 500
	in := &api.GetEventsInput{
		Streams: streamNames,
		Prev:    &limit,
	}
	for {
		output, err := workspace.GetEvents(ctx, in)
		if err != nil {
			return err
		}

		for _, event := range output.Items {
			t, err := time.Parse(chrono.RFC3339NanoUTC, event.Timestamp)
			if err != nil {
				cmdutil.Warnf("invalid event timestamp: %q", event.Timestamp)
				continue
			}

			var label string
			if showName {
				label = event.Stream
				if componentName := streamToLabel[event.Stream]; componentName != "" {
					label = componentName
				}
			}

			el.LogEvent(event.Stream, t, label, event.Message)
		}
		in.Cursor = &output.NextCursor
		in.Prev = nil

		if stopOnError {
			descriptions, err = workspace.DescribeProcesses(ctx, &api.DescribeProcessesInput{})
			if err != nil {
				return fmt.Errorf("failed to check status of processes: %w", err)
			}

			for _, proc := range descriptions.Processes {
				for _, id := range streamNames {
					if proc.ID == id && !proc.Running {
						return fmt.Errorf("process stopped running: %q", proc.Name)
					}
				}
			}
		}

		if len(output.Items) < 10 { // TODO: OK heuristic?
			if logFlags.NoFollow {
				return nil
			}
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				return nil
			}
		}
	}
}
