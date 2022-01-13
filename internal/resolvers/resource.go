package resolvers

import "context"

type ResourceResolver struct {
	Q *QueryResolver
	ResourceRow
}

type ResourceRow struct {
	IRI         string  `db:"iri"`
	ComponentID *string `db:"component_id"`
}

func (r *QueryResolver) AllResources(ctx context.Context) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT iri, component_id
		FROM resource
		ORDER BY iri ASC
	`)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ResourceResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ResourceResolver{
			Q:           r,
			ResourceRow: row,
		}
	}
	return resolvers, nil
}

func (r *QueryResolver) ResourceByIRI(ctx context.Context, args struct {
	IRI string
}) (*ResourceResolver, error) {
	return r.stackByIRI(ctx, &args.IRI)
}

func (r *QueryResolver) stackByIRI(ctx context.Context, iri *string) (*ResourceResolver, error) {
	s := &ResourceResolver{}
	err := r.getRowByID(ctx, &s.ResourceRow, `
		SELECT iri, component_id
		FROM resource
		WHERE iri = ?
	`, iri)
	if s.IRI == "" {
		s = nil
	}
	return s, err
}

func (r *ResourceResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByID(ctx, r.ComponentID)
}

func (r *QueryResolver) resourcesByStack(ctx context.Context, stackID string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT iri, component_id
		FROM resource
		WHERE stack_id = ?
		ORDER BY iri ASC
	`, stackID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ResourceResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ResourceResolver{
			Q:           r,
			ResourceRow: row,
		}
	}
	return resolvers, nil
}

func (r *QueryResolver) resourcesByComponent(ctx context.Context, componentID string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT iri, component_id
		FROM resource
		WHERE component_id = ?
		ORDER BY iri ASC
	`, componentID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ResourceResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ResourceResolver{
			Q:           r,
			ResourceRow: row,
		}
	}
	return resolvers, nil
}

func (r *QueryResolver) resourcesByProject(ctx context.Context, projectID string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT iri, component_id
		FROM resource
		INNER JOIN component ON component_id = component.id
		INNER JOIN stack ON component.stack_id = stack.id
		WHERE project_id = ?
		ORDER BY iri ASC
	`, projectID)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ResourceResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ResourceResolver{
			Q:           r,
			ResourceRow: row,
		}
	}
	return resolvers, nil
}
