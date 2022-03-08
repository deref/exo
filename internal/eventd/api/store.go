// Generated file. DO NOT EDIT.

package api

import (
	"context"
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
	Stream    string            `json:"stream"`
	Timestamp string            `json:"timestamp"`
	Message   string            `json:"message"`
	Tags      map[string]string `json:"tags"`
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

type StreamDescription struct {
	Name        string  `json:"name"`
	LastEventAt *string `json:"lastEventAt"`
}

type Event struct {
	ID        string            `json:"id"`
	Stream    string            `json:"stream"`
	Timestamp string            `json:"timestamp"`
	Message   string            `json:"message"`
	Tags      map[string]string `json:"tags"`
}
