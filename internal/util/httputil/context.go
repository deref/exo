package httputil

import (
	"context"
	"net/http"
	"time"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/telemetry"
	"github.com/deref/exo/internal/util/logging"
)

// TODO: Refactor this in to multiple middlewares.
func HandlerWithContext(ctx context.Context, handler http.Handler) http.Handler {
	debugHTTP := false // TODO: Configurable.
	logger := logging.CurrentLogger(ctx)
	tel := telemetry.FromContext(ctx)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := gensym.RandomBase32()
		sl := logger.Sublogger("http " + requestID)
		if debugHTTP {
			sl.Infof("%s %s", req.Method, req.URL)
		}
		ctx := logging.ContextWithLogger(req.Context(), sl)
		ctx = telemetry.ContextWithTelemetry(ctx, tel)
		start := chrono.Now(ctx)
		logw := &responseLogger{rw: w}
		handler.ServeHTTP(logw, req.WithContext(ctx))
		end := chrono.Now(ctx)
		duration := end.Sub(start).Truncate(time.Millisecond)
		if debugHTTP {
			sl.Infof("status %d - %s", logw.status, duration)
		}
	})
}

type responseLogger struct {
	rw     http.ResponseWriter
	status int
	size   int
}

func (rl *responseLogger) Header() http.Header {
	return rl.rw.Header()
}

func (rl *responseLogger) Write(bytes []byte) (int, error) {
	if rl.status == 0 {
		rl.status = http.StatusOK
	}
	size, err := rl.rw.Write(bytes)
	rl.size += size
	return size, err
}

func (rl *responseLogger) WriteHeader(status int) {
	rl.status = status
	rl.rw.WriteHeader(status)
}

func (rl *responseLogger) Flush() {
	if f, ok := rl.rw.(http.Flusher); ok {
		f.Flush()
	}
}
