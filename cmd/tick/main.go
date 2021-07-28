// Useful program for testing.

package main

import (
	"fmt"
	"time"
)

func main() {
	i := 0
	for {
		i++
		fmt.Printf("tick %d at %v\n", i, time.Now())
		<-time.After(1 * time.Second)
	}
}
