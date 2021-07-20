package main

import (
	"os"

	"github.com/deref/exo/fifofum"
)

func main() {
	command := os.Args[0]
	args := os.Args[1:]
	fifofum.Main(command, args)
}
