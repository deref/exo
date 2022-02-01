package resolvers

import (
	"context"
	"fmt"
	"sync"

	"github.com/deref/exo/internal/util/logging"
)

// Logger that creates event entities associated with a task.
type TaskLogger struct {
	Root      *RootResolver
	SystemLog logging.Logger
	TaskID    string

	// streamID is resolved lazily to avoid creating a majority of empty streams.
	mx       sync.RWMutex
	streamID string
}

var _ logging.Logger = (*TaskLogger)(nil)

func (tl *TaskLogger) Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	// XXX Remove this after watchJob and tailLogs are unified.
	tl.SystemLog.Infof("XXX TEMPORARY LOGGING task=%s: %s", tl.TaskID, message)
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
		Message:  message,
	})
	return err
}

func (tl *TaskLogger) resolveStream(ctx context.Context) (id string, err error) {
	// Get previously resolved stream id under a read-only lock.
	tl.mx.RLock()
	id = tl.streamID
	tl.mx.RUnlock()
	if id != "" {
		return
	}

	// Upgrade lock to read/write guard against write-race.
	tl.mx.Lock()
	defer tl.mx.Unlock()
	id = tl.streamID
	if id != "" {
		return
	}

	// Lazily create stream and cache result.
	var stream *StreamResolver
	stream, err = tl.Root.findOrCreateStream(ctx, "Task", tl.TaskID)
	if stream != nil {
		id = stream.ID
		tl.streamID = id
	}
	return
}

func (tl *TaskLogger) Sublogger(prefix string) logging.Logger {
	return &logging.Sublogger{
		Underlying: tl,
		Prefix:     prefix,
	}
}
