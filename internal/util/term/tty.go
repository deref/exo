package term

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsInteractive() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}
