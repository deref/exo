package resolvers

import (
	"context"
	"time"

	"github.com/deref/exo/internal/chrono"
	. "github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/logging"
)

func (r *QueryResolver) Now(ctx context.Context) Instant {
	return Now(ctx)
}

func (r *MutationResolver) Sleep(ctx context.Context, args struct {
	Seconds float64
}) (*VoidResolver, error) {
	// TODO: If sleeping over a certain amount of time, report "progress".
	// XXX
	return nil, chrono.Sleep(ctx, time.Duration(args.Seconds)*time.Second)
}

func (r *MutationResolver) Tick(ctx context.Context, args struct {
	Limit *int32
}) <-chan Instant {
	c := make(chan Instant)
	go func() {
		defer close(c)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for i := 0; ; i++ {
			if args.Limit != nil && int(*args.Limit) <= i {
				return
			}
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				logging.Infof(ctx, "tick %d at %s", i, t)
				select {
				case <-ctx.Done():
					return
				case c <- GoTimeToInstant(t):
				}
			}
		}
	}()
	return c
}
