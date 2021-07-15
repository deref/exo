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
	// Paginates events. Inputs before and after are mutually exclusive.
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
	Before string   `json:"before"`
	After  string   `json:"after"`
}

type GetEventsOutput struct {
	Events []Event `json:"events"`
	Cursor string  `json:"cursor"`
}

func NewLogCollectorMux(prefix string, iface LogCollector) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	BuildLogCollectorMux(b, iface)
	return b.Mux()
}

func BuildLogCollectorMux(b *josh.MuxBuilder, iface LogCollector) {
	b.AddMethod("add-log", iface.AddLog)
	b.AddMethod("remove-log", iface.RemoveLog)
	b.AddMethod("describe-logs", iface.DescribeLogs)
	b.AddMethod("get-events", iface.GetEvents)
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
