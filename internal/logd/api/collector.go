// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/internal/josh/server"
)

// Manages a set of logs. Collects and stores events from them.
type LogCollector interface {
	ClearEvents(context.Context, *ClearEventsInput) (*ClearEventsOutput, error)
	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	AddEvent(context.Context, *AddEventInput) (*AddEventOutput, error)
	// Returns pages of log events for some set of logs. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the log.
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
	RemoveOldEvents(context.Context, *RemoveOldEventsInput) (*RemoveOldEventsOutput, error)
}

type ClearEventsInput struct {
	Logs []string `json:"logs"`
}

type ClearEventsOutput struct {
}

type DescribeLogsInput struct {
	Names []string `json:"names"`
}

type DescribeLogsOutput struct {
	Logs []LogDescription `json:"logs"`
}

type AddEventInput struct {
	Log       string `json:"log"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type AddEventOutput struct {
}

type GetEventsInput struct {
	Logs   []string `json:"logs"`
	Cursor *string  `json:"cursor"`
	Prev   *int     `json:"prev"`
	Next   *int     `json:"next"`
}

type GetEventsOutput struct {
	Items      []Event `json:"items"`
	PrevCursor string  `json:"prevCursor"`
	NextCursor string  `json:"nextCursor"`
}

type RemoveOldEventsInput struct {
}

type RemoveOldEventsOutput struct {
}

func BuildLogCollectorMux(b *josh.MuxBuilder, factory func(req *http.Request) LogCollector) {
	b.AddMethod("clear-events", func(req *http.Request) interface{} {
		return factory(req).ClearEvents
	})
	b.AddMethod("describe-logs", func(req *http.Request) interface{} {
		return factory(req).DescribeLogs
	})
	b.AddMethod("add-event", func(req *http.Request) interface{} {
		return factory(req).AddEvent
	})
	b.AddMethod("get-events", func(req *http.Request) interface{} {
		return factory(req).GetEvents
	})
	b.AddMethod("remove-old-events", func(req *http.Request) interface{} {
		return factory(req).RemoveOldEvents
	})
}

type LogDescription struct {
	Name        string  `json:"name"`
	LastEventAt *string `json:"lastEventAt"`
}

type Event struct {
	ID        string `json:"id"`
	Log       string `json:"log"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
