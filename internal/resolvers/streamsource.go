package resolvers

import (
	"context"
	"fmt"
)

type StreamSourceResolver interface {
	Stream() *StreamResolver
	eventPrototype(ctx context.Context) (EventRow, error)
}

func (r *QueryResolver) findEventSource(ctx context.Context, typ string, id string) (StreamSourceResolver, error) {
	entity, err := r.findEntity(ctx, typ, id)
	if entity == nil || err != nil {
		return nil, err
	}
	underlying := entity.Underlying
	source, ok := underlying.(StreamSourceResolver)
	if !ok {
		return nil, fmt.Errorf("not a stream source: %T", underlying)
	}
	return source, nil
}
