// Generated file. DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	josh "github.com/deref/exo/internal/josh/server"
)

// Database of events organized into streams.
type Store interface {
	ClearEvents(context.Context, *ClearEventsInput) (*ClearEventsOutput, error)
	DescribeStreams(context.Context, *DescribeStreamsInput) (*DescribeStreamsOutput, error)
	AddEvent(context.Context, *AddEventInput) (*AddEventOutput, error)
	// Returns pages of events for some set of streams. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the stream.
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
	RemoveOldEvents(context.Context, *RemoveOldEventsInput) (*RemoveOldEventsOutput, error)
}

type ClearEventsInput struct {
	Streams []string `json:"streams"`
}

type ClearEventsOutput struct {
}

type DescribeStreamsInput struct {
	Names []string `json:"names"`
}

type DescribeStreamsOutput struct {
	Streams []StreamDescription `json:"streams"`
}

type AddEventInput struct {
	Stream    string `json:"stream"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type AddEventOutput struct {
}

type GetEventsInput struct {
	Streams   []string `json:"streams"`
	Cursor    *string  `json:"cursor"`
	FilterStr string   `json:"filterStr"`
	Prev      *int     `json:"prev"`
	Next      *int     `json:"next"`
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

func BuildStoreMux(b *josh.MuxBuilder, factory func(req *http.Request) Store) {
	b.AddMethod("clear-events", func(req *http.Request) interface{} {
		return factory(req).ClearEvents
	})
	b.AddMethod("describe-streams", func(req *http.Request) interface{} {
		return factory(req).DescribeStreams
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

type StreamDescription struct {
	Name        string  `json:"name"`
	LastEventAt *string `json:"lastEventAt"`
}

type Event struct {
	ID        string `json:"id"`
	Stream    string `json:"stream"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
