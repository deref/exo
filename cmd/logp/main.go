package main

import (
	"github.com/deref/exo/logrot"
	"github.com/deref/pier"
)

func main() {
	svc := logrot.NewService()
	pier.Main(logrot.NewHandler(svc))
}
