package resolvers

import (
	"context"
	"fmt"
	"math/rand"
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
	err := r.sleep(ctx, args.Seconds)
	return nil, err
}

func (r *MutationResolver) sleep(ctx context.Context, seconds float64) error {
	start := chrono.Now(ctx)
	end := start.Add(time.Duration(seconds) * time.Second)
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
			return nil
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

func (r *MutationResolver) BusyWork(ctx context.Context, args struct {
	Size     int32
	Width    *int32
	Depth    *int32
	Length   *int32
	FailRate *float64
}) (*VoidResolver, error) {
	size := int(args.Size)

	width := size
	if args.Width != nil {
		width = int(*args.Width)
	}

	depth := size
	if args.Depth != nil {
		depth = int(*args.Depth)
	}

	length := size
	if args.Length != nil {
		length = int(*args.Length)
	}

	failRate := 0.0
	if args.FailRate != nil {
		failRate = *args.FailRate
	}

	failRoll := rand.Float64()
	if failRoll < failRate {
		return nil, fmt.Errorf("random failure: %v", failRoll)
	}

	logging.Infof(ctx, "size=%d width=%d depth=%d length=%d", size, width, depth, length)

	for i := 0; i < width; i++ {
		d := rand.Intn(depth + 1)
		var err error
		if d > 0 {
			w := rand.Intn(width + 1)
			_, err = r.createTask(ctx, "busyWork", map[string]interface{}{
				"size":     args.Size,
				"width":    w,
				"depth":    d,
				"length":   length,
				"failRate": failRate,
			})
		}
		if err != nil {
			return nil, fmt.Errorf("creating subtask: %w", err)
		}
	}

	err := r.sleep(ctx, rand.Float64()*float64(length))
	return nil, err
}
