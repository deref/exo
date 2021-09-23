// Badger's default logger directly writes to directly to stderr.
// This file contains a badger.Logger implementation that wraps a
// Go builtin *log.Logger.

package badger

import (
	"github.com/deref/exo/internal/util/logging"
	"github.com/dgraph-io/badger/v3"
)

type logger struct {
	underlying logging.Logger
	level      int
}

const defaultLogLevel = int(badger.INFO)

func newLogger(underlying logging.Logger, level int) badger.Logger {
	return &logger{
		underlying: underlying,
		level:      level,
	}
}

func (l *logger) Printf(f string, v ...interface{}) {
	l.underlying.Infof("badger "+f, v...)
}

func (l *logger) Errorf(f string, v ...interface{}) {
	if l.level <= int(badger.ERROR) {
		l.Printf("ERROR: "+f, v...)
	}
}

func (l *logger) Warningf(f string, v ...interface{}) {
	if l.level <= int(badger.WARNING) {
		l.Printf("WARNING: "+f, v...)
	}
}

func (l *logger) Infof(f string, v ...interface{}) {
	if l.level <= int(badger.INFO) {
		l.Printf("INFO: "+f, v...)
	}
}

func (l *logger) Debugf(f string, v ...interface{}) {
	if l.level <= int(badger.DEBUG) {
		l.Printf("DEBUG: "+f, v...)
	}
}
