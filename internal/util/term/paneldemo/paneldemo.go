package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/deref/exo/internal/util/term"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	panel := &term.BottomPanel{}
	defer panel.Close()

	panel.SetHeight(8)

	animateHeight := true

	<-time.After(1 * time.Second)
	tick := 0
	for {
		tick++
		now := time.Now()
		fmt.Printf("tick %d at %v\n", tick, now)

		if animateHeight {
			h := int(math.Abs(math.Sin(float64(tick)/16.0)) * 10)
			panel.SetHeight(h)
		}

		var sb strings.Builder
		for line := 0; line < 4; line++ {
			fmt.Fprintf(&sb, "content line %d at time: %v\n", line, now)
		}
		panel.SetContent(sb.String())

		select {
		case <-time.After(100 * time.Millisecond):
		case <-ctx.Done():
			return
		}
	}
}
