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
}

var _ logging.Logger = (*EventLogger)(nil)

func (el *EventLogger) Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if err := el.createEvent(message); err != nil {
		el.SystemLog.Infof("error logging event for %s %s: %v", el.SourceType, el.SourceID, err)
		el.SystemLog.Infof("attempted message was: %s", message)
	}
}

func (el *EventLogger) createEvent(message string) error {
	// Avoid infinite regress, should Service perform any logging.
	ctx := context.Background()
	ctx = logging.ContextWithLogger(ctx, el.SystemLog)

	var createEvent struct {
		Event struct {
			ID string
		} `graphql:"createEvent(sourceType: $sourceType, sourceID: $sourceID, type: $type, message: $message)"`
	}
	return Mutate(ctx, el.Service, &createEvent, map[string]interface{}{
		"sourceType": el.SourceType,
		"sourceID":   el.SourceID,
		"type":       "Message",
		"message":    message,
	})
}

func (sl *EventLogger) Sublogger(prefix string) logging.Logger {
	return &logging.Sublogger{
		Underlying: sl,
		Prefix:     prefix,
	}
}
