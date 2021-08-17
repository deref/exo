// Useful program for testing.

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/deref/exo/internal/util/cmdutil"
)

func main() {
	cmd, err := cmdutil.ParseArgs(os.Args)
	if err != nil {
		cmdutil.Fatalf("parsing arguments: %w", err)
	}

	timeoutMS := 1000
	if timeoutFlag, ok := cmd.Flags["timeout-ms"]; ok {
		var err error
		if timeoutMS, err = strconv.Atoi(timeoutFlag); err != nil {
			cmdutil.Fatalf("parsing --timeout-ms: %w", err)
		}
	}

	i := 0
	for {
		i++
		fmt.Printf("tick %d at %v\n", i, time.Now())
		<-time.After(time.Duration(timeoutMS) * time.Millisecond)
	}
}
