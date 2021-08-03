package logging

import (
	golog "log"
)

type Logger interface {
	Infof(format string, v ...interface{})
}

type GoLogger struct {
	Underlying *golog.Logger
}

func (l *GoLogger) Infof(format string, v ...interface{}) {
	l.Underlying.Printf(format, v...)
}

func Default() Logger {
	return &GoLogger{
		Underlying: golog.Default(),
	}
}
