package cli

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"time"

	"github.com/Nerdmaster/terminal"
	"github.com/aybabtme/rgbterm"
	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/core/api"
	joshclient "github.com/deref/exo/internal/core/client"
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/term"
	"github.com/lucasb-eyer/go-colorful"
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

	colors := NewColorCache()

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
	labelWidth := 0
	streamToLabel := make(map[string]string, len(descriptions.Processes))
	streamToLabel[workspaceID] = "EXO"
	for _, process := range descriptions.Processes {
		streamToLabel[process.ID] = process.Name
		if labelWidth < len(process.Name) {
			labelWidth = len(process.Name)
		}
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
			timestamp := t.Local().Format("15:04:05")

			var prefix string
			if showName {
				label := event.Stream
				if componentName := streamToLabel[event.Stream]; componentName != "" {
					label = componentName
				} else if labelWidth < len(label) {
					labelWidth = len(label)
				}
				label = fmt.Sprintf("%*s", labelWidth, label)
				color := colors.Color(event.Stream)
				r, g, b := color.RGB255()
				prefix = rgbterm.FgString(
					fmt.Sprintf("%s %s", timestamp, label),
					r, g, b,
				)
			} else {
				prefix = timestamp
			}

			fmt.Printf("%s %s%s\r\n", prefix, event.Message, termReset)
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

type ColorCache struct {
	palette []colorful.Color
	colors  map[string]colorful.Color
}

func NewColorCache() *ColorCache {
	pal, err := colorful.HappyPalette(256)
	if err != nil {
		// An error should only be possible if the number of colours requested is
		// too high. Since this is a fixed constant this panic should be impossible.
		panic(err)
	}
	return &ColorCache{
		palette: pal,
		colors:  make(map[string]colorful.Color),
	}
}

func (cache *ColorCache) Color(key string) colorful.Color {
	color := cache.colors[key]
	if colorIsBlack(color) {
		b := md5.Sum([]byte(key))[0]
		color = cache.palette[b]
		cache.colors[key] = color
	}
	return color
}

func colorIsBlack(c colorful.Color) bool {
	return c.R == 0 && c.G == 0 && c.B == 0
}

const termReset = "\u001b[0m"
