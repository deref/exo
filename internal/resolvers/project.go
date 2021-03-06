package resolvers

import (
	"context"
	"fmt"

	"github.com/deref/exo/internal/gensym"
)

type ProjectResolver struct {
	Q *QueryResolver
	ProjectRow
}

type ProjectRow struct {
	ID          string  `db:"id"`
	DisplayName *string `db:"display_name"`
}

func projectRowsToResolvers(r *RootResolver, rows []ProjectRow) []*ProjectResolver {
	resolvers := make([]*ProjectResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ProjectResolver{
			Q:          r,
			ProjectRow: row,
		}
	}
	return resolvers
}

func (r *QueryResolver) AllProjects(ctx context.Context) ([]*ProjectResolver, error) {
	var rows []ProjectRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM project
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	return projectRowsToResolvers(r, rows), nil
}

func (r *QueryResolver) ProjectByID(ctx context.Context, args struct {
	ID string
}) (*ProjectResolver, error) {
	return r.projectByID(ctx, &args.ID)
}

func (r *QueryResolver) projectByID(ctx context.Context, id *string) (*ProjectResolver, error) {
	proj := &ProjectResolver{
		Q: r,
	}
	err := r.getRowByKey(ctx, &proj.ProjectRow, `
		SELECT *
		FROM project
		WHERE id = ?
	`, id)
	if proj.ID == "" {
		proj = nil
	}
	return proj, err
}

func (r *QueryResolver) projectByRef(ctx context.Context, ref string) (*ProjectResolver, error) {
	workspace, err := r.workspaceByRef(ctx, &ref)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, nil
	}
	return workspace.Project(ctx)
}

func (r *MutationResolver) CreateProject(ctx context.Context, args struct {
	DisplayName *string
}) (*ProjectResolver, error) {
	var row ProjectRow
	row.ID = gensym.RandomBase32()
	row.DisplayName = trimmedPtr(args.DisplayName, row.ID)
	if err := r.insertRow(ctx, "project", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &ProjectResolver{
		Q:          r,
		ProjectRow: row,
	}, nil
}

func (r *ProjectResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	return r, nil
}

func (r *ProjectResolver) Stacks(ctx context.Context) ([]*StackResolver, error) {
	return r.Q.stacksByProject(ctx, r.ID)
}

func (r *ProjectResolver) stackByRef(ctx context.Context, ref string) (*StackResolver, error) {
	return r.Q.stackByProjectIDAndRef(ctx, r.ID, ref)
}

func (r *ProjectResolver) Resources(ctx context.Context) ([]*ResourceResolver, error) {
	return r.Q.resourcesByProject(ctx, r.ID)
}
