package main

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/deref/exo/chrono"
	"github.com/deref/exo/exod/api"
	"github.com/deref/exo/util/cmdutil"
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

		colors := NewColorCache()

		logRefs := args
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
		logToComponent := make(map[string]string, len(descriptions.Processes))
		for _, process := range descriptions.Processes {
			// SEE NOTE: [LOG_COMPONENTS].
			for _, role := range []string{"out", "err"} {
				logName := fmt.Sprintf("%s:%s", process.ID, role)
				logToComponent[logName] = process.Name
			}
		}

		cursor := ""
		for {
			output, err := workspace.GetEvents(ctx, &api.GetEventsInput{
				Logs:  logIDs,
				After: cursor,
			})
			if err != nil {
				return err
			}

			for _, event := range output.Events {
				t, err := time.Parse(chrono.RFC3339NanoUTC, event.Timestamp)
				if err != nil {
					cmdutil.Warnf("invalid event timestamp: %q", event.Timestamp)
					continue
				}
				timestamp := t.Format("15:04:05")

				var prefix string
				if showName {
					label := event.Log
					if componentName := logToComponent[event.Log]; componentName != "" {
						label = componentName
					}
					color := colors.Color(label)
					prefix = rgbterm.FgString(
						fmt.Sprintf("%s %s", timestamp, label),
						color.Red, color.Green, color.Blue,
					)
				} else {
					prefix = timestamp
				}

				fmt.Printf("%s %s\n", prefix, event.Message)
			}
			cursor = output.Cursor
			if len(output.Events) < 10 { // TODO: OK heuristic?
				<-time.After(250 * time.Millisecond)
			}
		}
	},
}

type ColorCache struct {
	pallet []Color
	colors map[string]Color
}

func NewColorCache() *ColorCache {
	return &ColorCache{
		pallet: makePallet(),
		colors: make(map[string]Color),
	}
}

func makePallet() []Color {
	n := 256
	pallet := make([]Color, n)
	for i := 0; i < n; i++ {
		h := float64(i) / float64(n)
		r, g, b := rgbterm.HSLtoRGB(h, 0.7, 0.5)
		pallet[i] = Color{r, g, b}
	}
	return pallet
}

func (cache *ColorCache) Color(key string) Color {
	color := cache.colors[key]
	if color.IsBlack() {
		b := md5.Sum([]byte(key))[0]
		color = cache.pallet[b]
		cache.colors[key] = color
	}
	return color
}

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (c Color) IsBlack() bool {
	return c.Red == 0 && c.Green == 0 && c.Blue == 0
}
