package main

import (
	"context"

	"github.com/deref/exo/core"
	"github.com/deref/pier"
)

func main() {
	ctx := core.NewContext(context.Background())
	pier.Main(core.NewHandler(ctx))
}
