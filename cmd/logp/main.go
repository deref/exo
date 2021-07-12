package main

import (
	"github.com/deref/exo/logrot"
	"github.com/deref/pier"
)

func main() {
	pier.Main(logrot.NewHandler())
}
