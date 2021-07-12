// TODO: Generate this with JOSH tools.

package logrot

import (
	"context"
	"net/http"

	"github.com/deref/exo/josh"
)

type Service interface {
	AddLog(context.Context, *AddLogInput) (*AddLogOutput, error)
	RemoveLog(context.Context, *RemoveLogInput) (*RemoveLogOutput, error)
	DescribeLogs(context.Context, *DescribeLogsInput) (*DescribeLogsOutput, error)
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
}

type AddLogInput struct {
	ID         string `json:"id"`
	SourcePath string `json:"sourcePath"`
}

type AddLogOutput struct{}

type RemoveLogInput struct {
	ID string `json:"id"`
}

type RemoveLogOutput struct{}

type DescribeLogsInput struct {
	IDs []string `json:"ids"`
}

type DescribeLogsOutput struct {
	Logs []LogDescription `json:"logs"`
}

type LogDescription struct {
	ID          string  `json:"id"`
	SourcePath  string  `json:"sourcePath"`
	LastEventAt *string `json:"lastEventAt"`
}

type GetEventsInput struct {
	LogIDs []string `json:"logIds"`
	Before string   `json:"before"`
	After  string   `json:"after"`
}

type GetEventsOutput struct {
	Events []Event `json:"events"`
}

type Event struct {
	SID       string `json:"sid"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func NewMux(prefix string, service Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(prefix+"add-log", josh.NewMethodHandler(service.AddLog))
	mux.Handle(prefix+"remove-log", josh.NewMethodHandler(service.RemoveLog))
	mux.Handle(prefix+"describe-logs", josh.NewMethodHandler(service.DescribeLogs))
	mux.Handle(prefix+"get-events", josh.NewMethodHandler(service.GetEvents))
	return mux
}
