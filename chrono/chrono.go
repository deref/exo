package chrono

import (
	"context"
	"time"
)

func NowString(ctx context.Context) string {
	return time.Now().Format("2006-01-02T15:04:05.999999999Z")
}
