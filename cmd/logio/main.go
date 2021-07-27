package main

import (
	"os"

	"github.com/deref/exo/logio"
)

func main() {
	command := os.Args[0]
	args := os.Args[1:]
	logio.Main(command, args)
}
