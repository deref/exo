package cli

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/deref/exo/internal/util/logging"
)

// Wraps underlying logger, teeing output to stderr when in debug mode.
// Supports suppressing this tee'd output when some region of code has
// requested exclusive terminal access via beginExclusive/endExclusive.
type Logger struct {
	underlying logging.Logger
}

var exclusionDepth int32 = 0

func beginExclusive() {
	atomic.AddInt32(&exclusionDepth, 1)
}

func endExclusive() {
	atomic.AddInt32(&exclusionDepth, 1)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if logToStderr() && atomic.LoadInt32(&exclusionDepth) == 0 {
		fmt.Fprintf(os.Stderr, format+"\n", v...)
	}
	l.underlying.Infof(format, v...)
}

func (l *Logger) Sublogger(prefix string) logging.Logger {
	return &Logger{
		underlying: l.underlying.Sublogger(prefix),
	}
}
