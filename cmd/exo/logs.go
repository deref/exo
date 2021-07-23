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
				color := colors.Color(event.Log)

				t, err := time.Parse(chrono.RFC3339NanoUTC, event.Timestamp)
				if err != nil {
					cmdutil.Warnf("invalid event timestamp: %q", event.Timestamp)
					continue
				}
				timestamp := t.Format("15:04:05")

				prefix := fmt.Sprintf("%s %s", timestamp, event.Log)
				fmt.Printf("%s %s\n",
					rgbterm.FgString(prefix, color.Red, color.Green, color.Blue),
					event.Message,
				)
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
