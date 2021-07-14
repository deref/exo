// For testing non-worker bits of logd.

package main

import (
	"context"

	"github.com/deref/exo/logcol/api"
	"github.com/deref/exo/logcol/server"
	"github.com/deref/pier"
)

func main() {
	ctx := context.Background()
	collector := server.NewLogCollector()
	collector.Start(ctx)
	defer collector.Stop(ctx)
	pier.Main(api.NewLogCollectorMux("/", collector))
}
