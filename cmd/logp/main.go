package main

import (
	"context"

	"github.com/deref/exo"
	"github.com/deref/pier"
)

func main() {
	ctx := exo.NewContext(context.Background())
	pier.Main(logrot.Handler())
}
