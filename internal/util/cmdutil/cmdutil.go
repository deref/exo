package cmdutil

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/config"
)

func Warnf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf(format, v...))
}

func Fatalf(format string, v ...interface{}) {
	Warnf(format, v...)
	os.Exit(1)
}

func Fatal(err error) {
	Fatalf("%v", err)
}

func GetAddr(cfg *config.Config) string {
	return fmt.Sprintf("localhost:%d", cfg.HTTPPort)
}

func MustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
