package chrono

import (
	"context"
	"time"
)

func Now(ctx context.Context) time.Time {
	return time.Now()
}

func NowNano(ctx context.Context) uint64 {
	return uint64(Now(ctx).UnixNano())
}

func NowString(ctx context.Context) string {
	return IsoNano(Now(ctx))
}

func NanoToIso(nano int64) string {
	return IsoNano(time.Unix(0, nano))
}

func IsoNano(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.999999999Z")
}
