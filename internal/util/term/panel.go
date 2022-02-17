package term

import (
	"fmt"
	"strings"

	"github.com/deref/exo/internal/util/mathutil"
)

// Zero value is ready to use. Operations on panel affect Stdout.
// Must call Close() when done, to restore terminal settings.
type BottomPanel struct {
	content string
	// Desired panel height.
	height int
	// Current terminal bottom margin height.
	margin int
	// Current terminal size.
	termHeight int
}

func (p *BottomPanel) Height() int {
	return p.margin
}

// Set the panel height. The rendered height may be smaller than the requested
// height, if the requested height exceeds the terminal height minus ScrollOff.
func (p *BottomPanel) SetHeight(value int) {
	for p.height < value {
		fmt.Println()
		p.height++
	}
	p.height = value
	p.render()
}

// Minimum number of scrolling terminal lines to display, no matter the
// requested panel size.
const ScrollOff = 4

func (p *BottomPanel) Content() string {
	return p.content
}

// Set the content to appear at the bottom of the scrollback buffer.
// Will be truncated if it is larger than the panel's height.
func (p *BottomPanel) SetContent(value string) {
	p.content = value
	p.render()
}

func (p *BottomPanel) render() {
	termWidth, termHeight := GetSize()

	fmt.Printf("%c[0J", Esc)

	// Set margin.
	oldMargin := p.margin
	margin := mathutil.IntMin(p.height, termHeight-ScrollOff)
	bottom := termHeight - margin
	if margin != p.margin || termHeight != p.termHeight {
		x := bottom
		if p.height == 0 {
			x = 0
		}
		fmt.Printf("%c[0;%dr", Esc, x)
		p.margin = margin
		p.termHeight = termHeight
	}

	// Fill content with trailing blank lines.
	lines := strings.Split(p.content, "\n")
	for i := 0; i < margin; i++ {
		line := ""
		if i < len(lines) {
			line += TrimToVisualLength(lines[i], termWidth)
		}
		// Move cursor, write string, then erase to end of line.
		fmt.Printf("%c[%d;1H%s%c[0K", Esc, bottom+i+1, line, Esc)
	}

	// Move cursor to end of scrollback.
	oldBottom := termHeight - oldMargin
	fmt.Printf("%c[%d;1H", Esc, oldBottom)
}

func (p *BottomPanel) Close() {
	p.height = 0
	p.content = ""
	p.render()
}
