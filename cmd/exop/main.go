package main

import (
	"context"

	"github.com/deref/exo/kernel"
	"github.com/deref/pier"
)

func main() {
	ctx := kernel.NewContext(context.Background())
	pier.Main(kernel.NewHandler(ctx))
}
