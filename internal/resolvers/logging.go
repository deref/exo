package resolvers

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/util/logging"
)

// Logger that creates event entities associated with a task.
type TaskLogger struct {
	Root      *RootResolver
	SystemLog logging.Logger
	TaskID    string
}

var _ logging.Logger = (*TaskLogger)(nil)

func (tl *TaskLogger) Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	if err := tl.createEvent(message); err != nil {
		tl.SystemLog.Infof("error logging event for task %s: %w", tl.TaskID, err)
		tl.SystemLog.Infof("attempted message was: %s", message)
	}
}

func (tl *TaskLogger) createEvent(message string) error {
	ctx := context.Background()
	ctx = logging.ContextWithLogger(ctx, tl.SystemLog)

	streamID, err := tl.resolveStream(ctx)
	if err != nil {
		return fmt.Errorf("resolving stream: %w", err)
	}
	_, err = tl.Root.CreateEvent(ctx, createEventArgs{
		StreamID: streamID,
		Type:     "Message",
		Message:  message,
	})
	return err
}

func (tl *TaskLogger) Sublogger(prefix string) logging.Logger {
	return &logging.Sublogger{
		Underlying: tl,
		Prefix:     prefix,
	}
}
