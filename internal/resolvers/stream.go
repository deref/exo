package resolvers

import (
	"context"
	"fmt"
)

type StreamResolver struct {
	Q          *RootResolver
	SourceType string `db:"source_type"`
	SourceID   string `db:"source_id"`
}

func (r *QueryResolver) streamForSource(typ string, id string) *StreamResolver {
	return &StreamResolver{
		Q:          r,
		SourceType: typ,
		SourceID:   id,
	}
}

func (r *StreamResolver) Source(ctx context.Context) (source StreamSourceResolver, err error) {
	var entity *Entity
	entity, err = r.Q.findEntity(ctx, r.SourceType, r.SourceID)
	if entity != nil {
		source = entity.Underlying.(StreamSourceResolver)
	}
	return
}

// TODO: Consider delegating to source?
func (r *StreamResolver) eventFilter() eventFilter {
	var res eventFilter
	var p *string
	switch r.SourceType {
	case "Workspace":
		p = &res.WorkspaceID
	case "Stack":
		p = &res.StackID
	case "Component":
		p = &res.ComponentID
	case "Job":
		p = &res.JobID
	case "Task":
		p = &res.TaskID
	default:
		panic(fmt.Errorf("unexpected stream source type: %q", r.SourceType))
	}
	*p = r.SourceID
	return res
}

func (r *StreamResolver) Message(ctx context.Context) (string, error) {
	event, err := r.Q.latestEvent(ctx, r.eventFilter())
	if event == nil || err != nil {
		return "", err
	}
	return event.Message, nil
}

func (r *StreamResolver) Events(ctx context.Context, args struct {
	Cursor    *string
	Prev      *int32
	Next      *int32
	IContains *string
}) (*EventPageResolver, error) {
	q := eventQuery{
		Filter: r.eventFilter(),
	}
	if args.Cursor != nil {
		q.Cursor = *args.Cursor
	}
	if args.Prev != nil {
		q.Prev = int(*args.Prev)
	}
	if args.Next != nil {
		q.Next = int(*args.Next)
	}
	if args.IContains != nil {
		q.Filter.IContains = *args.IContains
	}
	return r.Q.findEvents(ctx, q)
}
