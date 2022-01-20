package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/util/jsonutil"
)

type ResourceResolver struct {
	Q *QueryResolver
	ResourceRow
}

type ResourceRow struct {
	ID        string  `db:"id"`
	Type      string  `db:"type"`
	OwnerType *string `db:"owner_type"`
	OwnerID   *string `db:"owner_id"`
	TaskID    *string `db:"task_id"`
	Model     string  `db:"model"`
	Status    *int    `db:"status"`
	Message   *string `db:"message"`
}

func (r *MutationResolver) ForgetResource(ctx context.Context, args struct {
	ID string
}) (*VoidResolver, error) {
	_, err := r.DB.ExecContext(ctx, `
		DELETE FROM resource
		WHERE id = ?
	`, args.ID)
	return nil, err
}

func (r *QueryResolver) AllResources(ctx context.Context) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT
			id,
			type,
			owner_type,
			owner_id,
			task_id,
			model,
			status,
			message
		FROM resource
		ORDER BY id ASC
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
	return r.resourceById(ctx, &args.IRI)
}

func (r *QueryResolver) resourceById(ctx context.Context, id *string) (*ResourceResolver, error) {
	s := &ResourceResolver{}
	err := r.getRowByID(ctx, &s.ResourceRow, `
		SELECT
			id,
			type,
			owner_type,
			owner_id,
			task_id,
			model,
			status,
			message
		FROM resource
		WHERE = ?
	`, id)
	if s.ID == "" {
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
		SELECT
			id,
			type,
			owner_type,
			owner_id,
			task_id,
			model,
			status,
			message
		FROM resource
		WHERE stack_id = ?
		ORDER BY id ASC
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
		SELECT
			id,
			type,
			owner_type,
			owner_id,
			task_id,
			model,
			status,
			message
		FROM resource
		WHERE component_id = ?
		ORDER BY id ASC
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
		SELECT
			resource.id,
			resource.type,
			resource.owner_type,
			resource.owner_id,
			resource.task_id,
			resource.model,
			resource.status,
			resource.message
		FROM resource
		INNER JOIN component ON component_id = component.id
		INNER JOIN stack ON component.stack_id = stack.id
		WHERE resource.project_id = ?
		ORDER BY resource.id ASC
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

func (r *MutationResolver) NewResource(ctx context.Context, args struct {
	Type      string
	Model     string
	OwnerType *string
	Workspace *string
	Component *string
	Adopt     *bool
}) (*ResourceResolver, error) {
	var row ResourceRow
	row.ID = gensym.RandomBase32()
	row.Type = args.Type

	// Lock the resource by pre-assigning a task ID.
	taskID := newTaskID()
	row.TaskID = &taskID

	adopt := args.Adopt != nil && *args.Adopt
	if adopt {
		row.Model = args.Model
	}

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
			return nil, fmt.Errorf("resolving project: %w", err)
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
		INSERT INTO resource (
			id,
			type,
			owner_type,
			owner_id,
			task_id,
			model,
			status,
			message
		) VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )
	`, row.ID, row.Type, row.OwnerType, row.OwnerID, row.TaskID, row.Model, row.Status, row.Message); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}

	parentTaskID := (*string)(nil)
	var err error
	if adopt {
		if _, err := r.newTask(ctx, taskID, parentTaskID, "readExternalResource", jsonutil.MustMarshalString(map[string]interface{}{
			"internalId": row.ID,
		})); err != nil {
			r.Logger.Infof("error starting resource %s adoption: %w", row.ID, err)
		}
	} else {
		if _, err = r.newTask(ctx, taskID, parentTaskID, "createExternalResource", jsonutil.MustMarshalString(map[string]interface{}{
			"internalId": row.ID,
			"model":      args.Model,
		})); err != nil {
			r.Logger.Infof("error starting resource %s creation: %w", row.ID, err)
		}
	}

	return &ResourceResolver{
		Q:           r,
		ResourceRow: row,
	}, nil
}

func (r *MutationResolver) lockResource(ctx context.Context, resourceID string) (unlock func(), err error) {
	currentTaskID := CurrentTaskID(ctx)
	if currentTaskID == nil {
		return nil, errors.New("synchronous mutations cannot lock resources")
	}
	taskID := *currentTaskID

	// Acquire or confirm previous acquisition of resource lock.
	res, err := r.DB.ExecContext(ctx, `
		UPDATE resource
		SET task_id = ?
		WHERE id = ?
		AND task_id IS NULL OR task_id = ?
	`, taskID, resourceID, taskID)
	if err != nil {
		return nil, fmt.Errorf("updating resource lock: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if n != 1 {
		return nil, errors.New("unable to lock resource")
	}

	return func() {
		if _, err := r.DB.ExecContext(ctx, `
			UPDATE resource
			SET task_id = NULL
			WHERE task_id = ?
		`, taskID); err != nil {
			r.Logger.Infof("task %s failed to unlock resource %s: %w", taskID, resourceID, err)
		}
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

func (r *MutationResolver) DisposeResource(ctx context.Context, args struct {
	IRI string
}) (*VoidResolver, error) {
	return nil, errors.New("TODO: implement DisposeResource")
}
