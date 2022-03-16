package resolvers

import (
	"context"
	"runtime"

	"github.com/deref/exo/internal/util/logging"
)

type GraphqlLogger struct {
	log logging.Logger
}

func NewGraphqlLogger(logger logging.Logger) *GraphqlLogger {
	return &GraphqlLogger{
		log: logger,
	}
}

// Adapted from graphql-gopher's DefaultLogger implementation.
func (gl *GraphqlLogger) LogPanic(ctx context.Context, value any) {
	const size = 64 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	gl.log.Infof("graphql: panic occurred: %v\n%s", value, buf)
}
