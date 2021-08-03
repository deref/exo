// Server unbundled from the CLI for development.
// The production entrypoint is the `exo server` command.

package main

import (
	"context"

	"github.com/deref/exo/exod"
	"github.com/deref/exo/util/logging"
)

func main() {
	ctx := context.Background()
	ctx = logging.ContextWithLogger(ctx, logging.Default())
	exod.Main(ctx)
}
