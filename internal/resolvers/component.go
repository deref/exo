package resolvers

import "context"

type ComponentResolver struct {
	Q *QueryResolver
	ComponentRow
}

type ComponentRow struct {
	ID      string `db:"id"`
	StackID string `db:"stack_id"`
	Name    string `db:"name"`
}

func (r *QueryResolver) ComponentByID(ctx context.Context, args struct {
	ID string
}) (*ComponentResolver, error) {
	return r.componentByID(ctx, &args.ID)
}

func (r *QueryResolver) componentByID(ctx context.Context, id *string) (*ComponentResolver, error) {
	s := &ComponentResolver{}
	err := r.getRowByID(ctx, &s.ComponentRow, `
		SELECT id, stack_id, name
		FROM component
		WHERE id = ?
	`, id)
	if s.ID == "" {
		s = nil
	}
	return s, err
}

func (r *ComponentResolver) Stack(ctx context.Context) (*StackResolver, error) {
	return r.Q.stackByID(ctx, &r.StackID)
}
