package main

import (
	"context"
	"net/http"

	"github.com/deref/exo/components/log"
	josh "github.com/deref/exo/josh/client"
	"github.com/deref/exo/kernel/server"
	logcol "github.com/deref/exo/logcol/client"
	"github.com/deref/pier"
)

func main() {
	ctx := server.NewContext(context.Background())
	ctx = log.ContextWithLogCollector(ctx, logcol.NewLogCollector(&josh.Client{
		HTTP: http.DefaultClient,
		URL:  "http://localhost:3001/",
	}))
	pier.Main(server.NewHandler(ctx))
}
