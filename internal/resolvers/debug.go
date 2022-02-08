package resolvers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/deref/exo/internal/api"
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

func (r *MutationResolver) BusyWork(ctx context.Context, args struct {
	Size   int32
	Width  *int32
	Depth  *int32
	Length *int32
}) (*VoidResolver, error) {
	ctxVars := api.CurrentContextVariables(ctx)
	var taskID *string
	if ctxVars.TaskID != "" {
		taskID = &ctxVars.TaskID
	}

	n := int(args.Size)

	width := n
	if args.Width != nil {
		width = int(*args.Width)
	}
	width = rand.Intn(width + 1)

	depth := n
	if args.Depth != nil {
		depth = int(*args.Depth)
	}

	length := n
	if args.Length != nil {
		length = int(*args.Length)
	}

	for i := 0; i < width; i++ {
		d := rand.Intn(depth + 1)
		var err error
		if d == 0 {
			_, err = r.createTask(ctx, "", taskID, "sleep", map[string]interface{}{
				"seconds": rand.Intn(length + 1),
			})
		} else {
			w := rand.Intn(width + 1)
			_, err = r.createTask(ctx, "", taskID, "busyWork", map[string]interface{}{
				"size":   args.Size,
				"width":  w,
				"depth":  d,
				"length": length,
			})
		}
		if err != nil {
			return nil, fmt.Errorf("creating subtask: %w", err)
		}
	}
	return nil, nil
}
