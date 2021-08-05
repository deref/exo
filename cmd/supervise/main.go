package main

import (
	"os"

	"github.com/deref/exo/internal/supervise"
)

func main() {
	command := os.Args[0]
	args := os.Args[1:]
	supervise.Main(command, args)
}
