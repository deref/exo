// TODO: Can we avoid putting the event store in to context?

package log

import (
	"context"

	"github.com/deref/exo/internal/eventd/api"
)

type contextKey int

const eventStoreKey contextKey = 1

func ContextWithEventStore(ctx context.Context, sto api.Store) context.Context {
	return context.WithValue(ctx, eventStoreKey, sto)
}

func CurrentEventStore(ctx context.Context) api.Store {
	return ctx.Value(eventStoreKey).(api.Store)
}
