package resolvers

import (
	"context"
	"fmt"
)

type Entity struct {
	Type       string
	ID         string
	Underlying interface{}
}

type EntityRef struct {
	Type string
	ID   string
}

func (r *QueryResolver) FindEntity(ctx context.Context, args struct {
	Type string
	ID   string
}) (*Entity, error) {
	return r.findEntity(ctx, args.Type, args.ID)
}

func (r *QueryResolver) findEntity(ctx context.Context, typ string, id string) (*Entity, error) {
	res := Entity{
		Type: typ,
		ID:   id,
	}
	var err error
	switch typ {
	case "System":
		res.Underlying = r.System()
	case "Task":
		res.Underlying, err = r.taskByID(ctx, &id)
	default:
		return nil, fmt.Errorf("cannot find entities of type: %q", typ)
	}
	if res.Underlying == nil || err != nil {
		return nil, err
	}
	return &res, nil
}
