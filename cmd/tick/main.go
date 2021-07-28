// Useful program for testing.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				fmt.Printf("Received signal %d\n", sig)
			}
		}
	}()

	i := 0
	for {
		i++
		fmt.Printf("tick %d at %v\n", i, time.Now())
		<-time.After(1 * time.Second)
	}
}
