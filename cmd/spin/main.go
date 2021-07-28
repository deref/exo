// Test program that pegs a core and ignores interrupts.

package main

import (
	"os"
	"os/signal"
)

func main() {
	signal.Ignore(os.Interrupt, os.Kill)
	for {
	}
}
