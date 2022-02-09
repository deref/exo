package api

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/util/logging"
)

// Logger that creates message events for a given source.
type EventLogger struct {
	Service    Service
	SystemLog  logging.Logger
	SourceType string
	SourceID   string
	Prefix     string
}

var _ logging.Logger = (*EventLogger)(nil)

func (el *EventLogger) Infof(format string, v ...interface{}) {
	format = logging.Prefix(el.Prefix, format)
	message := fmt.Sprintf(format, v...)
	if err := el.createEvent(message); err != nil {
		panic(err)
	}
}

func (el *EventLogger) createEvent(message string) error {
	ctx := context.Background()
	ctx = logging.ContextWithLogger(ctx, el.SystemLog)

	var createEvent struct {
		Event struct {
			ID string
		} `graphql:"createEvent(sourceType: $sourceType, sourceId: $sourceId, type: $type, message: $message)"`
	}
	return Mutate(ctx, el.Service, &createEvent, map[string]interface{}{
		"sourceType": el.SourceType,
		"sourceId":   el.SourceID,
		"type":       "Message",
		"message":    message,
	})
}

func (sl *EventLogger) Sublogger(prefix string) logging.Logger {
	sub := *sl
	sub.Prefix = logging.Prefix(sl.Prefix, prefix)
	return &sub
}

func NewSystemLogger(svc Service) *EventLogger {
	return &EventLogger{
		Service:    svc,
		SystemLog:  &logging.NopLogger{}, // Avoid infinite regress.
		SourceType: "System",
		SourceID:   "SYSTEM",
	}
}
