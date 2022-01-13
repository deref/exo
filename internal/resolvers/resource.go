package resolvers

import (
	"context"
	"errors"
	"fmt"
)

type ResourceResolver struct {
	Q *QueryResolver
	ResourceRow
}

type ResourceRow struct {
	IRI       string  `db:"iri"`
	OwnerType *string `db:"owner_type"`
	OwnerID   *string `db:"owner_id"`
}

func (r *QueryResolver) AllResources(ctx context.Context) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT iri, owner_type, owner_id
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
		SELECT iri, owner_type, owner_id
		FROM resource
		WHERE iri = ?
	`, iri)
	if s.IRI == "" {
		s = nil
	}
	return s, err
}

func (r *ResourceResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	if r.OwnerType == nil || *r.OwnerType != "component" {
		return nil, nil
	}
	return r.Q.componentByID(ctx, r.OwnerID)
}

func (r *QueryResolver) resourcesByStack(ctx context.Context, stackID string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT iri, owner_type, owner_id
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
		SELECT iri, owner_type, owner_id
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
		SELECT iri, owner_type, owner_id
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

func (r *MutationResolver) AdoptResource(ctx context.Context, args struct {
	IRI       string
	Project   *string
	Stack     *string
	Component *string
}) (*ResourceResolver, error) {
	var row ResourceRow
	row.IRI = args.IRI

	var component *ComponentResolver
	if args.Component != nil {
		var err error
		component, err = r.componentByRef(ctx, *args.Component, args.Stack)
		if err != nil {
			return nil, fmt.Errorf("resolving component: %q", err)
		}
		if component == nil {
			return nil, fmt.Errorf("no such component: %q", *args.Component)
		}
		row.OwnerType = stringPtr("component")
		row.OwnerID = stringPtr(component.ID)
	}

	var stack *StackResolver
	if args.Stack != nil || component != nil {
		var err error
		if args.Stack != nil {
			stack, err = r.stackByRef(ctx, *args.Stack)
		} else if component != nil {
			stack, err = component.Stack(ctx)
		}
		if err != nil {
			return nil, fmt.Errorf("resolving stack: %w", err)
		}
		if stack == nil {
			return nil, errors.New("stack not found")
		}
		if component != nil && component.StackID != stack.ID {
			return nil, fmt.Errorf("component %q not part of stack %q", *args.Component, *args.Stack)
		}
		if row.OwnerType == nil {
			row.OwnerType = stringPtr("stack")
			row.OwnerID = stringPtr(stack.ID)
		}
	}

	var project *ProjectResolver
	if args.Project != nil || stack != nil {
		var err error
		if args.Project != nil {
			project, err = r.projectByRef(ctx, *args.Project)
		} else if stack != nil {
			project, err = stack.Project(ctx)
		}
		if err != nil {
			return nil, fmt.Errorf("resolving project: %w", err)
		}
		if project == nil {
			return nil, errors.New("project not found")
		}
		if stack != nil && (stack.ProjectID == nil || *stack.ProjectID != project.ID) {
			return nil, errors.New("stack does not belong to project")
		}
		if row.OwnerType == nil {
			row.OwnerType = stringPtr("project")
			row.OwnerID = stringPtr(project.ID)
		}
	}

	if _, err := r.DB.ExecContext(ctx, `
		INSERT INTO resource ( iri, owner_type, owner_id )
		VALUES ( ?, ?, ? )
	`, row.IRI, row.OwnerType, row.OwnerID); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}
	return &ResourceResolver{
		Q:           r,
		ResourceRow: row,
	}, nil
}

func (r *ResourceResolver) Owner(ctx context.Context) (interface{}, error) {
	if r.OwnerType == nil {
		return nil, nil
	}
	switch *r.OwnerType {
	case "component":
		return r.Q.componentByID(ctx, r.OwnerID)
	case "stack":
		return r.Q.stackByID(ctx, r.OwnerID)
	case "project":
		return r.Q.projectByID(ctx, r.OwnerID)
	default:
		return nil, fmt.Errorf("unexpected owner type: %q", *r.OwnerType)
	}
}

func (r *ResourceResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	owner, err := r.Owner(ctx)
	if owner == nil || err != nil {
		return nil, err
	}
	return owner.(interface {
		Project(ctx context.Context) (*ProjectResolver, error)
	}).Project(ctx)
}

func (r *ResourceResolver) Stack(ctx context.Context) (*StackResolver, error) {
	owner, err := r.Owner(ctx)
	if owner == nil || err != nil {
		return nil, err
	}
	return owner.(interface {
		Stack(ctx context.Context) (*StackResolver, error)
	}).Stack(ctx)
}
