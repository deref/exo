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
	if r.OwnerType == nil || *r.OwnerType != "Component" {
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
	OwnerType *string
	Workspace *string
	Component *string
}) (*ResourceResolver, error) {
	var row ResourceRow
	row.IRI = args.IRI

	var workspace *WorkspaceResolver
	if args.Workspace != nil {
		var err error
		workspace, err := r.workspaceByRef(ctx, *args.Workspace)
		if err != nil {
			return nil, fmt.Errorf("resolving workspace: %w", err)
		}
		if workspace == nil {
			return nil, errors.New("no such workspace")
		}
	}

	var component *ComponentResolver
	if args.Component != nil {
		if workspace == nil {
			return nil, errors.New("workspace is required if component is provided")
		}
		var err error
		component, err = workspace.componentByRef(ctx, *args.Component)
		if err != nil {
			return nil, fmt.Errorf("resolving component: %w", err)
		}
		if component == nil {
			return nil, errors.New("no such component")
		}
	}

	var stack *StackResolver
	if workspace != nil {
		var err error
		stack, err = workspace.Stack(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving stack: %w", err)
		}
	}

	var project *ProjectResolver
	if workspace != nil {
		var err error
		project, err = workspace.Project(ctx)
		if err != nil {
			return nil, fmt.Errorf("resolving stack: %w", err)
		}
	}

	effectiveOwnerType := ""
	if args.OwnerType == nil {
		if component != nil {
			effectiveOwnerType = "Component"
		} else if stack != nil {
			effectiveOwnerType = "Stack"
		} else if project != nil {
			effectiveOwnerType = "Project"
		}
	} else {
		effectiveOwnerType = *args.OwnerType
	}
	row.OwnerType = stringPtr(effectiveOwnerType)
	switch effectiveOwnerType {
	case "":
		row.OwnerType = nil
	case "Component":
		if component == nil {
			return nil, errors.New("no component to set owner to")
		}
		row.OwnerID = stringPtr(component.ID)
	case "Stack":
		if stack == nil {
			return nil, errors.New("no stack to set owner to")
		}
		row.OwnerID = stringPtr(stack.ID)
	case "Project":
		if project == nil {
			return nil, errors.New("no project to set owner to")
		}
		row.OwnerID = stringPtr(project.ID)
	default:
		return nil, fmt.Errorf("unexpected owner type: %q", *args.OwnerType)
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
	case "Component":
		return r.Q.componentByID(ctx, r.OwnerID)
	case "Stack":
		return r.Q.stackByID(ctx, r.OwnerID)
	case "Project":
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
