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
	start := chrono.Now(ctx)
	end := start.Add(time.Duration(args.Seconds) * time.Second)
	deadline, cancel := context.WithDeadline(ctx, end)
	progress := ProgressInput{
		Total: int32(end.Sub(start).Milliseconds()),
	}
	defer cancel()
	for {
		select {
		case <-deadline.Done():
			progress.Current = progress.Total
			r.reportProgress(ctx, progress)
			return nil, nil
		case <-time.After(250 * time.Millisecond):
			progress.Current = int32(chrono.Now(ctx).Sub(start).Milliseconds())
			r.reportProgress(ctx, progress)
		}
	}
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
