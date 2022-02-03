package resolvers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/gensym"
	"github.com/deref/exo/internal/providers/sdk"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/util-go/httputil"
)

type ResourceResolver struct {
	Q *QueryResolver
	ResourceRow
}

type ResourceRow struct {
	ID        string  `db:"id"`
	Type      string  `db:"type"`
	IRI       *string `db:"iri"`
	OwnerType *string `db:"owner_type"`
	OwnerID   *string `db:"owner_id"`
	JobID     *string `db:"job_id"`
	Model     string  `db:"model"`
	Status    int32   `db:"status"`
	Message   *string `db:"message"`
}

func (r *MutationResolver) ForgetResource(ctx context.Context, args struct {
	Ref string
}) (*VoidResolver, error) {
	return r.forgetResource(ctx, args.Ref)
}

func (r *MutationResolver) forgetResource(ctx context.Context, ref string) (*VoidResolver, error) {
	_, err := r.DB.ExecContext(ctx, `
		DELETE FROM resource
		WHERE id = ?
		OR iri = ?
	`, ref, ref)
	return nil, err
}

func (r *QueryResolver) AllResources(ctx context.Context) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
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

func (r *QueryResolver) resourcesByIRI(ctx context.Context, iri string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.DB.SelectContext(ctx, &rows, `
		SELECT *
		FROM resource
		WHERE iri = ?
		ORDER BY id ASC
	`, iri)
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

func (r *QueryResolver) ResourceByID(ctx context.Context, args struct {
	ID string
}) (*ResourceResolver, error) {
	return r.resourceByID(ctx, &args.ID)
}

func (r *QueryResolver) resourceByID(ctx context.Context, id *string) (*ResourceResolver, error) {
	s := &ResourceResolver{}
	err := r.getRowByKey(ctx, &s.ResourceRow, `
		SELECT *
		FROM resource
		WHERE id = ?
	`, id)
	if s.ID == "" {
		s = nil
	}
	return s, err
}

func (r *QueryResolver) ResourceByIRI(ctx context.Context, args struct {
	IRI string
}) (*ResourceResolver, error) {
	return r.resourceByIRI(ctx, &args.IRI)
}

func (r *QueryResolver) resourceByIRI(ctx context.Context, iri *string) (*ResourceResolver, error) {
	if iri == nil {
		return nil, nil
	}
	resources, err := r.resourcesByIRI(ctx, *iri)
	if err != nil {
		return nil, err
	}
	switch len(resources) {
	case 0:
		return nil, nil
	case 1:
		return resources[0], nil
	default:
		return nil, errors.New("ambiguous")
	}
}

func isIRI(s string) bool {
	return strings.Contains(s, ":")
}

func (r *QueryResolver) ResourceByRef(ctx context.Context, args struct {
	Ref string
}) (*ResourceResolver, error) {
	return r.resourceByRef(ctx, &args.Ref)
}

func (r *QueryResolver) resourceByRef(ctx context.Context, ref *string) (*ResourceResolver, error) {
	if ref == nil {
		return nil, nil
	}
	if strings.Contains(*ref, ":") {
		return r.resourceByIRI(ctx, ref)
	}
	return r.resourceByID(ctx, ref)
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
		SELECT *
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
		SELECT *
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
		SELECT resource.*
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

func (r *ResourceResolver) Job() *JobResolver {
	return r.Q.jobByID(r.JobID)
}

func (r *ResourceResolver) Operation(ctx context.Context) (*string, error) {
	task, err := r.Job().Task(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving task: %w", err)
	}
	if task == nil {
		return nil, nil
	}
	var operation string
	switch task.Mutation {
	case "initializeResource":
		operation = "creating"
	case "refreshResource":
		operation = "reading"
	case "updateResource":
		operation = "updating"
	case "deleteResource":
		operation = "deleting"
	default:
		operation = "unknown"
	}
	return &operation, err
}

func (r *MutationResolver) CreateResource(ctx context.Context, args struct {
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

	// Lock the resource by pre-assigning a job id.
	jobID := newTaskID()
	row.JobID = &jobID

	adopt := args.Adopt != nil && *args.Adopt
	if adopt {
		row.Model = args.Model
	}

	var workspace *WorkspaceResolver
	if args.Workspace != nil {
		var err error
		workspace, err := r.workspaceByRef(ctx, args.Workspace)
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

	if err := r.insertRow(ctx, "resource", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}

	var err error
	if adopt {
		if _, err := r.createJob(ctx, jobID, "refreshResource", map[string]interface{}{
			"ref": row.ID,
		}); err != nil {
			logging.Infof(ctx, "error starting resource %s adoption: %w", row.ID, err)
		}
	} else {
		if _, err = r.createJob(ctx, jobID, "initializeResource", map[string]interface{}{
			"ref":   row.ID,
			"model": args.Model,
		}); err != nil {
			logging.Infof(ctx, "error starting resource %s creation: %w", row.ID, err)
		}
	}

	return &ResourceResolver{
		Q:           r,
		ResourceRow: row,
	}, nil
}

type resourceOperation func(ctx context.Context, row *ResourceRow, controller *sdk.Controller) (model string, err error)

func (r *MutationResolver) InitializeResource(ctx context.Context, args struct {
	Ref   string
	Model string
}) (*ResourceResolver, error) {
	return r.doResourceOperation(ctx, args.Ref,
		func(ctx context.Context, row *ResourceRow, ctrl *sdk.Controller) (string, error) {
			return ctrl.Create(ctx, row.Model)
		},
	)
}

func (r *MutationResolver) RefreshResource(ctx context.Context, args struct {
	Ref string
}) (*ResourceResolver, error) {
	return r.doResourceOperation(ctx, args.Ref,
		func(ctx context.Context, row *ResourceRow, ctrl *sdk.Controller) (string, error) {
			return ctrl.Read(ctx, row.Model)
		},
	)
}

func (r *MutationResolver) UpdateResource(ctx context.Context, args struct {
	Ref   string
	Model string
}) (*ResourceResolver, error) {
	return r.doResourceOperation(ctx, args.Ref,
		func(ctx context.Context, row *ResourceRow, ctrl *sdk.Controller) (string, error) {
			return ctrl.Update(ctx, row.Model, args.Model)
		},
	)
}

func (r *MutationResolver) DisposeResource(ctx context.Context, args struct {
	Ref string
}) (*VoidResolver, error) {
	if _, err := r.doResourceOperation(ctx, args.Ref,
		func(ctx context.Context, row *ResourceRow, ctrl *sdk.Controller) (string, error) {
			err := ctrl.Delete(ctx, row.Model)
			return row.Model, err
		},
	); err != nil {
		return nil, err
	}
	return r.forgetResource(ctx, args.Ref)
}

func (r *MutationResolver) doResourceOperation(ctx context.Context, ref string, op resourceOperation) (_ *ResourceResolver, doErr error) {
	ctxVars := api.CurrentContextVariables(ctx)
	if ctxVars == nil || ctxVars.TaskID == "" {
		return nil, errors.New("resource operations must be asynchronous")
	}
	jobID := ctxVars.JobID
	if ctxVars.TaskID != jobID {
		// Sanity check to avoid re-entering lock.
		// TODO: Should resources by locked by tasks instead of jobs? Probably, so that these tasks can occur as part of reconciliation jobs.
		return nil, errors.New("resource operation tasks must be top-level job")
	}

	// Acquire resource lock or confirm prior acquisition.
	// For resource initialization, the lock will be pre-acquired, so that
	// createResource can return a resource ID synchronously.
	var row ResourceRow
	if err := r.DB.GetContext(ctx, &row, `
		UPDATE resource
		SET job_id = ?, status = 0, message = NULL
		WHERE (id = ? OR iri = ?)
		AND (job_id IS NULL OR job_id = ?)
		RETURNING *
	`, jobID, ref, ref, jobID); err != nil {
		return nil, fmt.Errorf("updating resource lock: %w", err)
	}

	finish := func() {
		var status int
		var message *string
		if doErr == nil {
			status = http.StatusOK
		} else {
			status = httputil.StatusOf(doErr)
			message = stringPtr(doErr.Error())
		}
		if _, err := r.DB.ExecContext(ctx, `
			UPDATE resource
			SET job_id = NULL, status = ?, message = ?
			WHERE job_id = ?
		`, status, message, jobID); err != nil {
			logging.Infof(ctx, "task %s failed to unlock resource %q: %w", jobID, ref, err)
		}
	}
	defer finish()

	ctrl := getController(ctx, row.Type)
	if ctrl == nil {
		return nil, fmt.Errorf("no controller for type: %q", row.Type)
	}

	model, err := op(ctx, &row, ctrl)
	if err != nil {
		return nil, fmt.Errorf("controller failed: %w", err)
	}

	iri, identifyErr := ctrl.Identify(ctx, model)
	iri = strings.TrimSpace(iri)
	if identifyErr == nil && iri != "" {
		if row.IRI != nil {
			iri = *row.IRI
			if row.IRI != nil && *row.IRI != iri {
				identifyErr = fmt.Errorf("cannot change IRI from %q to %q", *row.IRI, iri)
			}
		}
	}

	if err := r.DB.GetContext(ctx, &row, `
		UPDATE resource
		SET model = ?, iri = ?
		WHERE id = ?
		RETURNING *
	`, model, iri, row.ID); err != nil {
		return nil, fmt.Errorf("recording model: %w", err)
	}
	if identifyErr != nil {
		return nil, fmt.Errorf("identifying: %w", identifyErr)
	}

	return &ResourceResolver{
		Q:           r,
		ResourceRow: row,
	}, nil
}

func (r *MutationResolver) CancelResourceOperation(ctx context.Context, args struct {
	Ref string
}) (*ResourceResolver, error) {
	resource, err := r.resourceByRef(ctx, &args.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving resource: %w", err)
	}
	if resource == nil {
		return nil, errors.New("no such resource")
	}

	if resource.JobID != nil {
		if err := r.cancelJob(ctx, *resource.JobID); err != nil {
			return nil, fmt.Errorf("canceling job: %w", err)
		}

		if _, err := r.DB.ExecContext(ctx, `
			UPDATE resource
			SET job_id = NULL, status = 500, message = 'interrupted'
			WHERE id = ? OR iri = ?
		`, args.Ref, args.Ref); err != nil {
			return nil, fmt.Errorf("releasing resource lock: %w", err)
		}
		resource.JobID = nil
	}

	return resource, nil
}
