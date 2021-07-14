// TODO: Generate this with JOSH tools.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/josh/server"
)

// TODO: Bulk methods.

type LogCollector interface {
	AddLog(context.Context, *AddLogInput) (*AddLogOutput, error)
	RemoveLog(context.Context, *RemoveLogInput) (*RemoveLogOutput, error)
	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
	Collect(context.Context, *CollectInput) (*CollectOutput, error)
}

type AddLogInput struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type AddLogOutput struct{}

type RemoveLogInput struct {
	Name string `json:"name"`
}

type RemoveLogOutput struct{}

type DescribeLogsInput struct {
	Names []string `json:"names"`
}

type DescribeLogsOutput struct {
	Logs []LogDescription `json:"logs"`
}

type LogDescription struct {
	Name        string  `json:"name"`
	Source      string  `json:"source"`
	LastEventAt *string `json:"lastEventAt"`
}

type GetEventsInput struct {
	Logs   []string `json:"logs"`
	Before string   `json:"before"`
	After  string   `json:"after"`
}

type GetEventsOutput struct {
	Events []Event `json:"events"`
}

type Event struct {
	Log       string `json:"log"`
	SID       string `json:"sid"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type CollectInput struct{}

type CollectOutput struct{}

func NewLogCollectorMux(prefix string, collector LogCollector) *http.ServeMux {
	b := josh.NewMuxBuilder(prefix)
	BuildLogCollectorMux(b, collector)
	return b.Mux()
}

func BuildLogCollectorMux(b *josh.MuxBuilder, collector LogCollector) {
	b.AddMethod("add-log", collector.AddLog)
	b.AddMethod("remove-log", collector.RemoveLog)
	b.AddMethod("describe-logs", collector.DescribeLogs)
	b.AddMethod("get-events", collector.GetEvents)
	b.AddMethod("collect", collector.Collect)
}
