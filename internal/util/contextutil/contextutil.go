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

func (c noCancel) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// WithoutCancel creates a new context from `ctx` that is never cancelled
// and never times out.
func WithoutCancel(ctx context.Context) context.Context {
	return noCancel{ctx: ctx}
}
