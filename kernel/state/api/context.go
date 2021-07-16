package api

import "context"

type contextKey int

const storeKey contextKey = 1

func ContextWithStore(ctx context.Context, store Store) context.Context {
	return context.WithValue(ctx, storeKey, store)
}

func CurrentStore(ctx context.Context) Store {
	return ctx.Value(storeKey).(Store)
}
