package cmdutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/deref/exo/internal/config"
)

// Print a string with a trailing newline.
// Avoid printing an _extra_ newline.
func Show(s string) {
	if strings.HasSuffix(s, "\n") {
		fmt.Print(s)
	} else if s != "" {
		fmt.Println(s)
	}
}

func Warnf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v\n", fmt.Errorf(format, v...))
}

// TODO: Many usages of this should panic, not hard-exit.
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
