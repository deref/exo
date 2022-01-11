package resolvers

import "context"

type ProjectResolver struct {
	Q *QueryResolver
	ProjectRow
}

type ProjectRow struct {
	ID          string  `db:"id"`
	DisplayName *string `db:"display_name"`
}

func (r *QueryResolver) ProjectByID(ctx context.Context, args struct {
	ID string
}) (*ProjectResolver, error) {
	return r.projectByID(ctx, &args.ID)
}

func (r *QueryResolver) projectByID(ctx context.Context, id *string) (*ProjectResolver, error) {
	proj := &ProjectResolver{}
	err := r.getRowByID(ctx, &proj.ProjectRow, `
		SELECT id, display_name
		FROM project
		WHERE id = ?
	`, id)
	if proj.ID == "" {
		proj = nil
	}
	return proj, err
}
