package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/deref/exo/internal/util/mathutil"
	"github.com/deref/exo/internal/util/term"
)

type EventLogger struct {
	W io.Writer
	// Labels will be right-aligned within this width.  Will grow throughout
	// logging, but may be initialized to keep the initial output aligned.
	LabelWidth int

	colors *term.ColorCache
}

func (el *EventLogger) Init() {
	el.LabelWidth = mathutil.IntMax(el.LabelWidth, 10)

	if useColor() {
		el.colors = term.NewColorCache()
	}
}

func (el *EventLogger) LogEvent(sourceID string, t time.Time, label string, message string) {
	var prefix, suffix string
	timestamp := t.Local().Format("15:04:05")
	if label == "" {
		prefix = timestamp
	} else {
		el.LabelWidth = mathutil.IntMax(el.LabelWidth, len(label))
		label = fmt.Sprintf("%*s", el.LabelWidth, label)
		prefix = fmt.Sprintf("%s %s", timestamp, label)
		if el.colors != nil {
			r, g, b := el.colors.Color(sourceID).RGB255()
			prefix = rgbterm.FgString(prefix, r, g, b)
			suffix = term.ResetCode
		}
	}
	fmt.Fprintf(el.W, "%s %s%s\r\n", prefix, message, suffix)
}
