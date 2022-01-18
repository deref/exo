package chrono

import (
	"context"
	"time"
)

func Sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-time.After(duration):
		return nil
	case <-ctx.Done():
		return context.Canceled
	}
}
