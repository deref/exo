package resolvers

import "context"

type StackResolver struct {
	Q *QueryResolver
	StackRow
}

type StackRow struct {
	ID string `db:"id"`
}

func (r *QueryResolver) StackByID(ctx context.Context, args struct {
	ID string
}) (*StackResolver, error) {
	return r.stackByID(ctx, &args.ID)
}

func (r *QueryResolver) stackByID(ctx context.Context, id *string) (*StackResolver, error) {
	s := &StackResolver{}
	err := r.getRowByID(ctx, &s.StackRow, `
		SELECT id
		FROM stack
		WHERE id = ?
	`, id)
	if s.ID == "" {
		s = nil
	}
	return s, err
}
