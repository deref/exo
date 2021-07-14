package cmdutil

import (
	"fmt"
	"os"
)

func Fatalf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}
