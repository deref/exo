package main

import (
	"context"

	"github.com/deref/exo/components/log"
	"github.com/deref/exo/kernel/server"
	logcol "github.com/deref/exo/logcol/server"
	"github.com/deref/pier"
)

func main() {
	ctx := server.NewContext(context.Background())
	ctx = log.ContextWithLogCollector(ctx, logcol.NewLogCollector())
	pier.Main(server.NewHandler(ctx))
}
