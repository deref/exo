package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/deref/exo/internal/util/mathutil"
	"github.com/deref/exo/internal/util/term"
)

type EventWriter struct {
	W io.Writer
	// Labels will be right-aligned within this width.  Will grow throughout
	// logging, but may be initialized to keep the initial output aligned.
	LabelWidth int

	colors *term.ColorCache
}

func (w *EventWriter) Init() {
	w.LabelWidth = mathutil.IntMax(w.LabelWidth, 10)

	if useColor() {
		w.colors = term.NewColorCache()
	}
}

func (w *EventWriter) PrintEvent(sourceID string, t time.Time, label string, message string) {
	var prefix, suffix string
	timestamp := t.Local().Format("15:04:05")
	if label == "" {
		prefix = timestamp
	} else {
		w.LabelWidth = mathutil.IntMax(w.LabelWidth, len(label))
		label = fmt.Sprintf("%*s", w.LabelWidth, label)
		prefix = fmt.Sprintf("%s %s", timestamp, label)
		if w.colors != nil {
			r, g, b := w.colors.Color(sourceID).RGB255()
			prefix = rgbterm.FgString(prefix, r, g, b)
			suffix = term.ResetCode
		}
	}
	fmt.Fprintf(w.W, "%s %s%s\r\n", prefix, message, suffix)
}
