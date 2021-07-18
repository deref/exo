// For testing non-worker bits of logd.

package main

import (
	"context"

	"github.com/deref/exo/cmdutil"
	"github.com/deref/exo/logd/api"
	"github.com/deref/exo/logd/server"
	"github.com/deref/pier"
)

func main() {
	paths := cmdutil.MustMakeDirectories()
	cfg := &server.Config{
		VarDir: paths.VarDir,
	}
	ctx := context.Background()
	collector := server.NewLogCollector(ctx, cfg)
	if err := collector.Start(ctx); err != nil {
		cmdutil.Fatalf("starting collector: %w", err)
	}
	defer collector.Stop(ctx)
	pier.Main(api.NewLogCollectorMux("/", collector))
}
