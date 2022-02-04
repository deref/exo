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
	CallDepth  int
}

func (l *GoLogger) Infof(format string, v ...interface{}) {
	l.Underlying.Output(2+l.CallDepth, fmt.Sprintf(format, v...))
}

func (l *GoLogger) Sublogger(prefix string) Logger {
	return &Sublogger{
		Underlying: l,
		Prefix:     prefix,
		CallDepth:  1,
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
	CallDepth  int
}

func (l *Sublogger) Infof(format string, v ...interface{}) {
	l.Underlying.Infof(l.Prefix+": "+format, v...)
}

func (l *Sublogger) Sublogger(prefix string) Logger {
	return &Sublogger{
		Underlying: l,
		Prefix:     l.Prefix + " " + prefix,
		CallDepth:  l.CallDepth + 1,
	}
}

type NopLogger struct{}

func (nop *NopLogger) Infof(format string, v ...interface{}) {}

func (nop *NopLogger) Sublogger(prefix string) Logger {
	return nop
}
