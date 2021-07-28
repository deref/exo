// Test program that blocks forever and ignores interrupts.

package main

import (
	"os"
	"os/signal"
	"time"
)

func main() {
	signal.Ignore(os.Interrupt, os.Kill)
	time.Sleep(24 * 365 * 100 * time.Hour)
}
