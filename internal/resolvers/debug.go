package resolvers

import (
	"context"
	"time"

	"github.com/deref/exo/internal/chrono"
)

func (r *MutationResolver) Sleep(ctx context.Context, args struct {
	Seconds float64
}) (*VoidResolver, error) {
	return nil, chrono.Sleep(ctx, time.Duration(args.Seconds)*time.Second)
}
