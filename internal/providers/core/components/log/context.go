// TODO: Delete this file & pass LogCollector around like state.Store.

package log

import (
	"context"

	"github.com/deref/exo/internal/logd/api"
)

type contextKey int

const logCollectorKey contextKey = 1

func ContextWithLogCollector(ctx context.Context, collector api.LogCollector) context.Context {
	return context.WithValue(ctx, logCollectorKey, collector)
}

func CurrentLogCollector(ctx context.Context) api.LogCollector {
	return ctx.Value(logCollectorKey).(api.LogCollector)
}
