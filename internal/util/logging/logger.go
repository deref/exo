package logging

import (
	"fmt"
	golog "log"
)

type Logger interface {
	Infof(format string, v ...interface{})
	Sublogger(prefix string) Logger
}

type GoLogger struct {
	Underlying *golog.Logger
}

func (l *GoLogger) Infof(format string, v ...interface{}) {
	l.Underlying.Output(3, fmt.Sprintf(format, v...))
}

func (l *GoLogger) Sublogger(prefix string) Logger {
	return &Sublogger{
		Underlying: l,
		Prefix:     prefix,
	}
}

func Default() Logger {
	return &GoLogger{
		Underlying: golog.Default(),
	}
}

type Sublogger struct {
	Underlying Logger
	Prefix     string
}

func (l *Sublogger) Infof(format string, v ...interface{}) {
	l.Underlying.Infof(l.Prefix+": "+format, v...)
}

func (l *Sublogger) Sublogger(prefix string) Logger {
	return &Sublogger{
		Underlying: l,
		Prefix:     l.Prefix + " " + prefix,
	}
}
