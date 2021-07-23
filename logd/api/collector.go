// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

// Manages a set of logs. Collects and stores events from them.
type LogCollector interface {
	AddLog(context.Context, *AddLogInput) (*AddLogOutput, error)
	RemoveLog(context.Context, *RemoveLogInput) (*RemoveLogOutput, error)
	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	// Returns pages of log events for some set of logs. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the log.
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
}

type AddLogInput struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type AddLogOutput struct {
}

type RemoveLogInput struct {
	Name string `json:"name"`
}

type RemoveLogOutput struct {
}

type DescribeLogsInput struct {
	Names []string `json:"names"`
}

type DescribeLogsOutput struct {
	Logs []LogDescription `json:"logs"`
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

func BuildLogCollectorMux(b *josh.MuxBuilder, factory func(req *http.Request) LogCollector) {
	b.AddMethod("add-log", func(req *http.Request) interface{} {
		return factory(req).AddLog
	})
	b.AddMethod("remove-log", func(req *http.Request) interface{} {
		return factory(req).RemoveLog
	})
	b.AddMethod("describe-logs", func(req *http.Request) interface{} {
		return factory(req).DescribeLogs
	})
	b.AddMethod("get-events", func(req *http.Request) interface{} {
		return factory(req).GetEvents
	})
}

type LogDescription struct {
	Name        string  `json:"name"`
	Source      string  `json:"source"`
	LastEventAt *string `json:"lastEventAt"`
}

type Event struct {
	ID        string `json:"id"`
	Log       string `json:"log"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
