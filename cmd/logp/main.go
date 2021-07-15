// For testing non-worker bits of logd.

package main

import (
	"context"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/logcol/api"
	"github.com/deref/exo/logcol/server"
	"github.com/deref/pier"
)

func main() {
	cfg := &server.Config{
		VarDir: cmdutil.MustVarDir(),
	}
	ctx := context.Background()
	collector := server.NewLogCollector(cfg)
	collector.Start(ctx)
	defer collector.Stop(ctx)
	pier.Main(api.NewLogCollectorMux("/", collector))
}
