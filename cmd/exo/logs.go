package main

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
	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/term"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
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
		checkOrEnsureServer()
		cl := newClient()
		workspace := requireCurrentWorkspace(ctx, cl)
		stopOnError := false
		return tailLogs(ctx, workspace, args, stopOnError)
	},
}

func tailLogs(ctx context.Context, workspace api.Workspace, logRefs []string, stopOnError bool) error {
	ctx, cancel := context.WithCancel(ctx)
	var eg errgroup.Group
	eg.Go(func() error {
		return runTailLogsReader(ctx, cancel)
	})
	eg.Go(func() error {
		return runTailLogsWriter(ctx, workspace, logRefs, stopOnError)
	})
	return eg.Wait()
}

func runTailLogsReader(ctx context.Context, cancel func()) error {
	stdin := os.Stdin.Fd()
	if !isatty.IsTerminal(stdin) {
		return nil
	}

	raw := &term.RawMode{
		FD: stdin,
	}
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

func runTailLogsWriter(ctx context.Context, workspace api.Workspace, logRefs []string, stopOnError bool) error {
	colors := NewColorCache()

	showName := len(logRefs) != 1

	resolved, err := workspace.Resolve(ctx, &api.ResolveInput{
		Refs: logRefs,
	})
	if err != nil {
		return fmt.Errorf("resolving refs: %w", err)
	}
	var logIDs []string
	if len(logRefs) > 0 {
		logIDs = make([]string, 0, len(logRefs))
		for _, logID := range resolved.IDs {
			if logID != nil {
				logIDs = append(logIDs, *logID)
			}
		}
	}

	descriptions, err := workspace.DescribeProcesses(ctx, &api.DescribeProcessesInput{})
	if err != nil {
		return fmt.Errorf("describing processes: %w", err)
	}
	labelWidth := 0
	logToComponent := make(map[string]string, len(descriptions.Processes))
	for _, process := range descriptions.Processes {
		if len(logRefs) == 0 {
			logIDs = append(logIDs, process.ID)
		}
		logToComponent[process.ID] = process.Name
		if labelWidth < len(process.Name) {
			labelWidth = len(process.Name)
		}
	}

	limit := 500
	in := &api.GetEventsInput{
		Streams: logIDs,
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
				if componentName := logToComponent[event.Stream]; componentName != "" {
					label = componentName
				} else if labelWidth < len(label) {
					labelWidth = len(label)
				}
				label = fmt.Sprintf("%*s", labelWidth, label)
				color := colors.Color(label)
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
				for _, id := range logIDs {
					if proc.ID == id && !proc.Running {
						return fmt.Errorf("process stopped running: %q", proc.Name)
					}
				}
			}
		}

		if len(output.Items) < 10 { // TODO: OK heuristic?
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
