package contextutil

import (
	"context"
	"time"
)

type noCancel struct {
	ctx context.Context
}

func (c noCancel) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c noCancel) Done() <-chan struct{} {
	return nil
}

func (c noCancel) Err() error {
	return nil
}

func (c noCancel) Value(key any) any {
	return c.ctx.Value(key)
}

// WithoutCancel creates a new context from `ctx` that is never cancelled
// and never times out.
// XXX This is an anti-pattern. It is used for background work, but that
// should always be indirected through some kind of background worker
// where the context must be explicitly propegated.
func WithoutCancel(ctx context.Context) context.Context {
	if res, ok := ctx.(noCancel); ok {
		return res
	}
	return noCancel{ctx: ctx}
}
