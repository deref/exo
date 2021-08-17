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

	intervalMS := 1000
	if intervalFlag, ok := cmd.Flags["interval-ms"]; ok {
		var err error
		if intervalMS, err = strconv.Atoi(intervalFlag); err != nil {
			cmdutil.Fatalf("parsing --interval-ms: %w", err)
		}
	}

	i := 0
	for {
		i++
		fmt.Printf("tick %d at %v\n", i, time.Now())
		<-time.After(time.Duration(intervalMS) * time.Millisecond)
	}
}
