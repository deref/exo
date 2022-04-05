package resolvers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/gensym"
	. "github.com/deref/exo/internal/scalars"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/exo/sdk"
	"github.com/deref/util-go/httputil"
)

type ResourceResolver struct {
	Q *QueryResolver
	ResourceRow
}

type ResourceRow struct {
	ID          string  `db:"id"`
	Type        string  `db:"type"`
	IRI         *string `db:"iri"`
	ProjectID   *string `db:"project_id"`
	StackID     *string `db:"stack_id"`
	ComponentID *string `db:"component_id"`
	TaskID      *string `db:"task_id"`
	RawModel    RawJSON `db:"model"`
	Status      int32   `db:"status"`
	Message     *string `db:"message"`
}

func (r *MutationResolver) ForgetResource(ctx context.Context, args struct {
	Ref string
}) (*VoidResolver, error) {
	err := r.forgetResource(ctx, args.Ref)
	return nil, err
}

func (r *MutationResolver) forgetResource(ctx context.Context, ref string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM resource
		WHERE id = ?
		OR iri = ?
	`, ref, ref)
	return err
}

func resourceRowsToResolvers(r *RootResolver, rows []ResourceRow) []*ResourceResolver {
	resolvers := make([]*ResourceResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &ResourceResolver{
			Q:           r,
			ResourceRow: row,
		}
	}
	return resolvers
}

func (r *QueryResolver) AllResources(ctx context.Context) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM resource
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	return resourceRowsToResolvers(r, rows), nil
}

func (r *QueryResolver) resourcesByIRI(ctx context.Context, iri string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM resource
		WHERE iri = ?
		ORDER BY id ASC
	`, iri)
	if err != nil {
		return nil, err
	}
	return resourceRowsToResolvers(r, rows), nil
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
		return nil, errutil.NewHTTPError(http.StatusConflict, "ambiguous resource iri")
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

func (r *QueryResolver) resourcesByProject(ctx context.Context, projectID string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM resource
		WHERE project_id = ?
		ORDER BY id ASC
	`, projectID)
	if err != nil {
		return nil, err
	}
	return resourceRowsToResolvers(r, rows), nil
}

func (r *QueryResolver) resourcesByStack(ctx context.Context, stackID string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM resource
		WHERE stack_id = ?
		ORDER BY id ASC
	`, stackID)
	if err != nil {
		return nil, err
	}
	return resourceRowsToResolvers(r, rows), nil
}

func (r *QueryResolver) resourcesByComponent(ctx context.Context, componentID string) ([]*ResourceResolver, error) {
	var rows []ResourceRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM resource
		WHERE component_id = ?
		ORDER BY id ASC
	`, componentID)
	if err != nil {
		return nil, err
	}
	return resourceRowsToResolvers(r, rows), nil
}

func (r *ResourceResolver) Model() (obj JSONObject, err error) {
	err = json.Unmarshal(r.RawModel, &obj)
	return
}

func (r *ResourceResolver) OwnerType() *string {
	var ownerType string
	switch {
	case r.ComponentID != nil:
		ownerType = "Component"
	case r.StackID != nil:
		ownerType = "Stack"
	case r.ProjectID != nil:
		ownerType = "Project"
	default:
		return nil
	}
	return &ownerType
}

func (r *ResourceResolver) Owner(ctx context.Context) (any, error) {
	ownerType := r.OwnerType()
	if ownerType == nil {
		return nil, nil
	}
	switch *ownerType {
	case "Component":
		return r.Component(ctx)
	case "Stack":
		return r.Stack(ctx)
	case "Project":
		return r.Project(ctx)
	default:
		panic(fmt.Errorf("unexpected owner type: %q", *ownerType))
	}
}

func (r *ResourceResolver) Project(ctx context.Context) (*ProjectResolver, error) {
	return r.Q.projectByID(ctx, r.ProjectID)
}

func (r *ResourceResolver) Stack(ctx context.Context) (*StackResolver, error) {
	return r.Q.stackByID(ctx, r.StackID)
}

func (r *ResourceResolver) Component(ctx context.Context) (*ComponentResolver, error) {
	return r.Q.componentByResourceID(ctx, &r.ID)
}

func (r *ResourceResolver) Task(ctx context.Context) (*TaskResolver, error) {
	return r.Q.taskByID(ctx, r.TaskID)
}

func (r *ResourceResolver) Operation(ctx context.Context) (*string, error) {
	task, err := r.Task(ctx)
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
		operation = "unexpected mutation: " + task.Mutation
	}
	return &operation, err
}

func (r *MutationResolver) CreateResource(ctx context.Context, args struct {
	Type      string
	Model     JSONObject
	Project   *string
	Stack     *string
	Component *string
	Adopt     *bool
}) (*ResourceResolver, error) {
	var row ResourceRow
	row.ID = gensym.RandomBase32()
	row.Type = args.Type

	adopt := args.Adopt != nil && *args.Adopt
	if adopt {
		row.RawModel = jsonutil.MustMarshal(args.Model)
	}

	var project *ProjectResolver
	if args.Project != nil {
		var err error
		project, err = r.projectByRef(ctx, *args.Project)
		if err := validateResolve("project", *args.Project, project, err); err != nil {
			return nil, err
		}
		row.ProjectID = &project.ID
	}

	var stack *StackResolver
	if args.Stack != nil {
		var err error
		if project == nil {
			stack, err = r.stackByRef(ctx, args.Stack)
		} else {
			stack, err = project.stackByRef(ctx, *args.Stack)
		}
		if err := validateResolve("stack", *args.Stack, stack, err); err != nil {
			return nil, err
		}
		row.StackID = &stack.ID
		if project == nil {
			row.ProjectID = stack.ProjectID
		}
	}

	var component *ComponentResolver
	if args.Component != nil {
		var err error
		if stack == nil {
			component, err = r.componentByRef(ctx, *args.Component, nil)
		} else {
			component, err = stack.componentByRef(ctx, *args.Component)
		}
		if err := validateResolve("component", *args.Component, component, err); err != nil {
			return nil, err
		}
		row.ComponentID = &component.ID
		if stack == nil {
			row.StackID = &component.StackID
		}
	}

	if err := r.insertRow(ctx, "resource", row); err != nil {
		return nil, fmt.Errorf("inserting: %w", err)
	}

	var job *JobResolver
	if adopt {
		var err error
		job, err = r.createJob(ctx, "refreshResource", map[string]any{
			"ref": row.ID,
		})
		if err != nil {
			logging.Infof(ctx, "error starting resource %s adoption: %w", row.ID, err)
		}
	} else {
		var err error
		job, err = r.createJob(ctx, "initializeResource", map[string]any{
			"ref":   row.ID,
			"model": args.Model,
		})
		if err != nil {
			logging.Infof(ctx, "error starting resource %s creation: %w", row.ID, err)
		}
	}

	resource, err := r.lockResource(ctx, row.ID, job.ID)
	if err != nil {
		r.SystemLog.Infof("error establishing initial resource lock: %w", err)
		resource = &ResourceResolver{
			Q:           r,
			ResourceRow: row,
		}
	}

	return resource, nil
}

type controlResourceFunc func(ctx context.Context, controller sdk.AResourceController, cfg *sdk.ResourceConfig, model *json.RawMessage) error

func (r *MutationResolver) InitializeResource(ctx context.Context, args struct {
	Ref   string
	Model JSONObject
}) (*ResourceResolver, error) {
	return r.controlResource(ctx, args.Ref,
		func(ctx context.Context, ctrl sdk.AResourceController, cfg *sdk.ResourceConfig, model *json.RawMessage) error {
			*model = jsonutil.MustMarshal(args.Model)
			return ctrl.CreateResource(ctx, cfg, model)
		},
	)
}

func (r *MutationResolver) RefreshResource(ctx context.Context, args struct {
	Ref string
}) (*ResourceResolver, error) {
	return r.controlResource(ctx, args.Ref,
		func(ctx context.Context, ctrl sdk.AResourceController, cfg *sdk.ResourceConfig, model *json.RawMessage) error {
			return ctrl.ReadResource(ctx, cfg, model)
		},
	)
}

func (r *MutationResolver) UpdateResource(ctx context.Context, args struct {
	Ref   string
	Model JSONObject
}) (*ResourceResolver, error) {
	return r.controlResource(ctx, args.Ref,
		func(ctx context.Context, ctrl sdk.AResourceController, cfg *sdk.ResourceConfig, model *json.RawMessage) error {
			prev := *model
			return ctrl.UpdateResource(ctx, cfg, &prev, model)
		},
	)
}

func (r *MutationResolver) DisposeResource(ctx context.Context, args struct {
	Ref string
}) (*VoidResolver, error) {
	if _, err := r.controlResource(ctx, args.Ref,
		func(ctx context.Context, ctrl sdk.AResourceController, cfg *sdk.ResourceConfig, model *json.RawMessage) error {
			return ctrl.DeleteResource(ctx, cfg, model)
		},
	); err != nil {
		return nil, err
	}
	err := r.forgetResource(ctx, args.Ref)
	return nil, err
}

func (r *ResourceResolver) controller(ctx context.Context) (sdk.AResourceController, error) {
	controller := r.Q.resourceControllerByType(ctx, r.Type)
	if controller == nil {
		return nil, fmt.Errorf("no resource controller for type: %q", r.Type)
	}
	return controller, nil
}

func (r *MutationResolver) controlResource(ctx context.Context, ref string, f controlResourceFunc) (_ *ResourceResolver, doErr error) {
	ctxVars := api.CurrentContextVariables(ctx)
	if ctxVars == nil || ctxVars.TaskID == "" {
		return nil, errors.New("resource operations must be asynchronous")
	}

	resource, err := r.lockResource(ctx, ref, ctxVars.TaskID)
	if err != nil {
		return nil, fmt.Errorf("acquiring resource lock: %w", err)
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
		if _, err := r.db.ExecContext(ctx, `
			UPDATE resource
			SET task_id = NULL, status = ?, message = ?
			WHERE task_id = ?
		`, status, message, resource.TaskID); err != nil {
			logging.Infof(ctx, "task %s failed to unlock resource %q: %w", resource.TaskID, ref, err)
		}
	}
	defer finish()

	ctrl, err := resource.controller(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving controller: %w", err)
	}

	cfg := &sdk.ResourceConfig{
		ID:   resource.ID,
		Type: resource.Type,
		IRI:  resource.IRI,
	}

	model := json.RawMessage(resource.RawModel)

	fErr := f(ctx, ctrl, cfg, &model)
	if fErr != nil {
		return nil, fmt.Errorf("controller failed: %w", err)
	}

	// Attempt to identify resource, even if f failed.
	iri, identifyErr := ctrl.IdentifyResource(ctx, cfg, &model)
	iri = strings.TrimSpace(iri)
	if identifyErr == nil && iri != "" {
		if resource.IRI != nil {
			iri = *resource.IRI
			if resource.IRI != nil && *resource.IRI != iri {
				identifyErr = fmt.Errorf("cannot change IRI from %q to %q", *resource.IRI, iri)
			}
		}
	}

	// Update model and iri, regardless of controller errors.
	var row ResourceRow
	if err := r.db.GetContext(ctx, &row, `
		UPDATE resource
		SET model = ?, iri = ?
		WHERE id = ?
		RETURNING *
	`, model, iri, resource.ID); err != nil {
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

func (r *MutationResolver) lockResource(ctx context.Context, ref string, taskID string) (*ResourceResolver, error) {
	// Acquire resource lock or confirm prior acquisition.
	var row ResourceRow
	if err := r.db.GetContext(ctx, &row, `
		UPDATE resource
		SET task_id = ?, status = 0, message = NULL
		WHERE (id = ? OR iri = ?)
		AND (task_id IS NULL OR task_id = ?)
		RETURNING *
	`, taskID, ref, ref, taskID); err != nil {
		return nil, err
	}
	if row.ID == "" {
		return nil, errors.New("resource unavailable")
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

	if resource.TaskID != nil {
		if err := r.cancelTask(ctx, *resource.TaskID); err != nil {
			return nil, fmt.Errorf("canceling job: %w", err)
		}

		if _, err := r.db.ExecContext(ctx, `
			UPDATE resource
			SET job_id = NULL, status = 500, message = 'interrupted'
			WHERE id = ? OR iri = ?
		`, args.Ref, args.Ref); err != nil {
			return nil, fmt.Errorf("releasing resource lock: %w", err)
		}
		resource.TaskID = nil
	}

	return resource, nil
}
