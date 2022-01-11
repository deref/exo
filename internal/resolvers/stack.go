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

func (r *QueryResolver) stacksByWorkspace(ctx context.Context, workspaceID string) ([]*StackResolver, error) {
	var rows []StackRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT id, project_id, root
		FROM workspace
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*StackResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &StackResolver{
			Q:        r,
			StackRow: row,
		}
	}
	return resolvers, nil
}
