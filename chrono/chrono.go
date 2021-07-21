package chrono

import (
	"context"
	"time"
)

const (
	RFC3339NanoUTC  = "2006-01-02T15:04:05.999999999Z"
	RFC3339MicroUTC = "2006-01-02T15:04:05.999999Z"
)

func Now(ctx context.Context) time.Time {
	return time.Now()
}

func NowNano(ctx context.Context) int64 {
	return Now(ctx).UnixNano()
}

func NowString(ctx context.Context) string {
	return IsoNano(Now(ctx))
}

func NanoToIso(nano int64) string {
	return IsoNano(time.Unix(0, nano))
}

func ParseIsoToNano(iso string) (int64, error) {
	t, err := time.Parse(RFC3339NanoUTC, iso)
	if err != nil {
		return 0, err
	}

	return t.UnixNano(), nil
}

func IsoNano(t time.Time) string {
	return t.Format(RFC3339NanoUTC)
}
