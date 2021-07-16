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
	paths := cmdutil.MustMakeDirectories()
	cfg := &server.Config{
		VarDir: paths.VarDir,
	}
	ctx := context.Background()
	collector := server.NewLogCollector(ctx, cfg)
	collector.Start(ctx)
	defer collector.Stop(ctx)
	pier.Main(api.NewLogCollectorMux("/", collector))
}
