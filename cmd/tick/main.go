// Useful program for testing.

package main

import (
	"log"
	"time"
)

func main() {
	i := 0
	for {
		i++
		log.Println("tick", i)
		<-time.After(1 * time.Second)
	}
}
