package logging

import (
	"fmt"
	golog "log"
)

type Logger interface {
	Infof(format string, v ...any)
	Sublogger(prefix string) Logger
}

type GoLogger struct {
	Prefix     string
	Underlying *golog.Logger
	CallDepth  int
}

func Prefix(prefix, format string) string {
	if prefix == "" {
		return format
	}
	return prefix + ": " + format
}

func (l *GoLogger) Infof(format string, v ...any) {
	format = Prefix(l.Prefix, format)
	l.Underlying.Output(2+l.CallDepth, fmt.Sprintf(format, v...))
}

func (l *GoLogger) Sublogger(prefix string) Logger {
	return &GoLogger{
		Underlying: l.Underlying,
		Prefix:     Prefix(l.Prefix, prefix),
		CallDepth:  l.CallDepth,
	}
}

func Default() Logger {
	return &GoLogger{
		Underlying: golog.Default(),
	}
}

type NopLogger struct{}

func (nop *NopLogger) Infof(format string, v ...any) {}

func (nop *NopLogger) Sublogger(prefix string) Logger {
	return nop
}
