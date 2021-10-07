package server

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/core/api"
	state "github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/deps"
	"github.com/deref/exo/internal/esv"
	eventd "github.com/deref/exo/internal/eventd/api"
	"github.com/deref/exo/internal/gensym"
	josh "github.com/deref/exo/internal/josh/server"
	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/procfile"
	"github.com/deref/exo/internal/providers/core"
	"github.com/deref/exo/internal/providers/core/components/invalid"
	"github.com/deref/exo/internal/providers/core/components/log"
	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/providers/docker/components/container"
	"github.com/deref/exo/internal/providers/docker/components/network"
	"github.com/deref/exo/internal/providers/docker/components/volume"
	"github.com/deref/exo/internal/providers/unix/components/process"
	"github.com/deref/exo/internal/task"
	"github.com/deref/exo/internal/util/contextutil"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	"github.com/deref/exo/internal/util/pathutil"
	dockerclient "github.com/docker/docker/client"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

type Workspace struct {
	ID          string
	VarDir      string
	Store       state.Store
	SyslogPort  uint
	Logger      logging.Logger // TODO: Embed in context, so it can be annotated with request info.
	Docker      *dockerclient.Client
	TaskTracker *task.TaskTracker
	EsvClient   *esv.EsvClient
}

func (ws *Workspace) logEventf(ctx context.Context, format string, v ...interface{}) {
	eventStore := log.CurrentEventStore(ctx)
	_, err := eventStore.AddEvent(ctx, &eventd.AddEventInput{
		Stream:    ws.ID,
		Timestamp: chrono.NowString(ctx),
		Message:   fmt.Sprintf(format, v...),
	})
	if err != nil {
		ws.Logger.Infof("error adding workspace event: %v", err)
	}
}

func (ws *Workspace) Describe(ctx context.Context, input *api.DescribeInput) (*api.DescribeOutput, error) {
	description, err := ws.describe(ctx)
	if err != nil {
		return nil, err
	}
	return &api.DescribeOutput{
		Description: *description,
	}, nil
}

func (ws *Workspace) describe(ctx context.Context) (*api.WorkspaceDescription, error) {
	output, err := ws.Store.DescribeWorkspaces(ctx, &state.DescribeWorkspacesInput{
		IDs: []string{ws.ID},
	})
	if err != nil {
		return nil, err
	}
	if len(output.Workspaces) != 1 {
		return nil, fmt.Errorf("invalid workspace: %q", ws.ID)
	}
	desc := output.Workspaces[0]
	return &api.WorkspaceDescription{
		ID:          ws.ID,
		Root:        desc.Root,
		DisplayName: desc.DisplayName,
	}, nil
}

func (ws *Workspace) Destroy(ctx context.Context, input *api.DestroyInput) (*api.DestroyOutput, error) {
	job := ws.TaskTracker.StartTask(ctx, "destroying")
	ws.logEventf(ctx, "destroying workspace... %s", job.JobID())
	query := makeComponentQuery(withReversedDependencies)

	go func() {
		defer job.Finish()
		ws.goControlComponents(job, query, func(*api.ComponentDescription) interface{} {
			return &api.DestroyInput{}
		})
		if err := job.Wait(); err != nil {
			return
		}
		if _, err := ws.Store.RemoveWorkspace(ctx, &state.RemoveWorkspaceInput{
			ID: ws.ID,
		}); err != nil {
			job.Fail(fmt.Errorf("removing workspace from store: %w", err))
			return
		}
	}()
	return &api.DestroyOutput{
		JobID: job.JobID(),
	}, nil
}

func (ws *Workspace) Apply(ctx context.Context, input *api.ApplyInput) (*api.ApplyOutput, error) {
	description, err := ws.describe(ctx)
	if err != nil {
		return nil, fmt.Errorf("describing workspace: %w", err)
	}
	res := ws.loadManifest(description.Root, input)
	if res.Err != nil {
		return nil, res.Err
	}
	m := res.Manifest

	describeOutput, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	oldComponents := make(map[string]api.ComponentDescription, len(describeOutput.Components))
	for _, oldComponent := range describeOutput.Components {
		oldComponents[oldComponent.Name] = oldComponent
	}

	// TODO: Handle partial failures.
	job := ws.TaskTracker.StartTask(ctx, "applying")
	ws.logEventf(ctx, "applying manifest... %s", job.JobID())

	// The algorithm for applying a new manifest is as follows:
	// 1. Build dependency graph for the new manifest, allComponents.
	// 2. Create empty dependency graphs for deletions, deleteGraph, and for creations, createGraph.
	// 3. Iterate all components in the old graph, and for each component that does not exist in the new graph:
	//   3.1. Add a node to deleteGraph
	//   3.2. For each dependency of node, add an _inverted_ edge for each dependency of that node. This
	//        allows us to check whether deleting these nodes would leave the graph with unmet dependencies.
	// 4. Check whether deleteGraph has unmet dependencies. If so, return an error because applying deletions would
	//    leave the graph with unmet dependencies.
	// 5. Create a set to track components to update, updateSet.
	// 6. Walk the new components.
	//   6.1. If the component already exists, get this component and all components that depend on it. For each,
	//     6.1.1. Check if the component is in updateSet. If not, add a delete task to deleteGraph and a create task
	//            to createGraph, and add the component to updateSet.
	//   6.2. Otherwise, add a task to create the component to createGraph.
	// 7. Apply all deletes in topographic order.
	// 8. Apply all creates in topographic order.

	// 1.
	allComponents := deps.New()
	for _, c := range m.Components {
		allComponents.AddNode(&componentNode{
			component: c,
		})
		for _, dependency := range c.DependsOn {
			allComponents.AddEdge(c.Name, dependency)
		}
	}

	// 2.
	createGraph := deps.New()
	deleteGraph := deps.New()

	// 3.
	for name, oldComponent := range oldComponents {
		if allComponents.HasNode(name) {
			continue
		}

		deleteGraph.AddNode(&runTaskNode{
			name: name,
			task: job.CreateChild("deleting " + name),
			run: func(t *task.Task) error {
				return ws.control(job.Context, oldComponent, &api.DestroyInput{})
			},
		})
		// Invert the dependencies for deletions so that we can check whether there is anything
		// left in the graph that still depends on something that we want to delete.
		for _, dependency := range oldComponent.DependsOn {
			deleteGraph.AddEdge(dependency, name)
		}
	}

	// 4.
	unmetDeps := deleteGraph.UnmetDependencies()
	if len(unmetDeps) > 0 {
		return nil, fmt.Errorf("would remove components that are still depended on: %s", strings.Join(unmetDeps, ", "))
	}

	// 5.
	updateSet := make(map[string]struct{})

	recreateComponentOnce := func(name string, oldComponent api.ComponentDescription, newComponent manifest.Component) {
		// 6.1.1.
		if _, alreadyUpdated := updateSet[name]; alreadyUpdated {
			return
		}

		deleteGraph.AddNode(&runTaskNode{
			name: name,
			task: job.CreateChild("deleting " + name),
			run: func(t *task.Task) error {
				return ws.control(job.Context, oldComponent, &api.DestroyInput{})
			},
		})
		for _, dependency := range oldComponent.DependsOn {
			deleteGraph.AddEdge(dependency, name)
		}

		createGraph.AddNode(&runTaskNode{
			name: name,
			task: job.CreateChild("re-creating " + name),
			run: func(t *task.Task) error {
				// Should the replacement component get the old component's ID?
				_, err := ws.createComponent(t, newComponent, gensym.RandomBase32())
				return err
			},
		})
		for _, dependency := range newComponent.DependsOn {
			createGraph.AddEdge(name, dependency)
		}
		updateSet[name] = struct{}{}
	}

	// 6.
	for _, newComponent := range m.Components {
		newComponent := newComponent
		name := newComponent.Name

		if oldComponent, exists := oldComponents[name]; exists {
			name := name
			// 6.1.
			recreateComponentOnce(name, oldComponent, newComponent)

			componentsForUpdate := allComponents.Dependents(name)
			for updatedComponentName, componentForUpdate := range componentsForUpdate {
				forDelete := oldComponents[updatedComponentName]
				forCreate := componentForUpdate.(*componentNode).component
				recreateComponentOnce(updatedComponentName, forDelete, forCreate)
			}
		} else {
			// 6.2.
			createGraph.AddNode(&runTaskNode{
				name: name,
				task: job.CreateChild("adding " + name),
				run: func(t *task.Task) error {
					id := gensym.RandomBase32()
					_, err := ws.createComponent(t, newComponent, id)
					return err
				},
			})
			for _, dependency := range newComponent.DependsOn {
				createGraph.AddEdge(name, dependency)
			}
		}
	}

	// Run tasks.
	go func() {
		defer job.Finish()

		executeRunTasks(deleteGraph)
		executeRunTasks(createGraph)
	}()

	return &api.ApplyOutput{
		Warnings: res.Warnings,
		JobID:    job.ID(),
	}, nil
}

func executeRunTasks(g *deps.Graph) {
	layers := g.TopoSortedLayers()
	for _, layer := range layers {
		var wg sync.WaitGroup
		for _, node := range layer {
			runTask := node.(*runTaskNode)
			wg.Add(1)
			go func() {
				defer wg.Done()
				runTask.task.Start()
				defer runTask.task.Finish()
				if err := runTask.run(runTask.task); err != nil {
					runTask.task.Fail(err)
				}
			}()
		}
		wg.Wait()
	}
}

func (ws *Workspace) Resolve(ctx context.Context, input *api.ResolveInput) (*api.ResolveOutput, error) {
	storeOutput, err := ws.Store.Resolve(ctx, &state.ResolveInput{
		WorkspaceID: ws.ID,
		Refs:        input.Refs,
	})
	if err != nil {
		return nil, err
	}
	var output api.ResolveOutput
	output.IDs = make([]*string, len(storeOutput.IDs))
	for i, id := range storeOutput.IDs {
		output.IDs[i] = id
	}
	return &output, err
}

func (ws *Workspace) DescribeComponents(ctx context.Context, input *api.DescribeComponentsInput) (*api.DescribeComponentsOutput, error) {
	stateOutput, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID:         ws.ID,
		Refs:                input.Refs,
		Types:               input.Types,
		IncludeDependencies: input.IncludeDependencies,
		IncludeDependents:   input.IncludeDependents,
	})
	if err != nil {
		return nil, err
	}

	output := &api.DescribeComponentsOutput{
		Components: make([]api.ComponentDescription, len(stateOutput.Components)),
	}
	for i, component := range stateOutput.Components {
		output.Components[i] = api.ComponentDescription{
			ID:        component.ID,
			Name:      component.Name,
			Type:      component.Type,
			Spec:      component.Spec,
			State:     component.State,
			Created:   component.Created,
			DependsOn: component.DependsOn,
		}
	}
	return output, nil
}

// TODO: Argument should be type+id, no spec or state or anything like that.
func (ws *Workspace) newController(ctx context.Context, desc api.ComponentDescription) Controller {
	description, err := ws.describe(ctx)
	if err != nil {
		return &invalid.Invalid{
			Err: fmt.Errorf("workspace error: %w", err),
		}
	}
	env, err := ws.getEnvironment(ctx)
	if err != nil {
		return &invalid.Invalid{
			Err: fmt.Errorf("environment error: %w", err),
		}
	}
	base := core.ComponentBase{
		ComponentID:          desc.ID,
		ComponentName:        desc.Name,
		ComponentState:       desc.State,
		WorkspaceID:          ws.ID,
		WorkspaceRoot:        description.Root,
		WorkspaceEnvironment: env,
		Logger:               ws.Logger,
	}
	switch desc.Type {
	case "process":
		return &process.Process{
			ComponentBase: base,
			SyslogPort:    ws.SyslogPort,
		}

	case "container":
		return &container.Container{
			ComponentBase: docker.ComponentBase{
				ComponentBase: base,
				Docker:        ws.Docker,
			},
			SyslogPort: ws.SyslogPort,
		}

	case "network":
		return &network.Network{
			ComponentBase: docker.ComponentBase{
				ComponentBase: base,
				Docker:        ws.Docker,
			},
		}

	case "volume":
		return &volume.Volume{
			ComponentBase: docker.ComponentBase{
				ComponentBase: base,
				Docker:        ws.Docker,
			},
		}

	default:
		return &invalid.Invalid{
			Err: fmt.Errorf("unsupported component type: %q", desc.Type),
		}
	}
}

func (ws *Workspace) DescribeEnvironment(ctx context.Context, input *api.DescribeEnvironmentInput) (*api.DescribeEnvironmentOutput, error) {
	env, err := ws.getEnvironment(ctx)
	if err != nil {
		return nil, err
	}
	return &api.DescribeEnvironmentOutput{
		Variables: env,
	}, nil
}

func (ws *Workspace) getEnvironment(ctx context.Context) (map[string]string, error) {
	envPath, err := ws.resolveWorkspacePath(ctx, ".env")
	if err != nil {
		return nil, fmt.Errorf("resolving env file path: %w", err)
	}

	env := map[string]string{}

	if ws.EsvClient != nil {
		// TODO: store this in the state.
		secretConfigPath, err := ws.resolveWorkspacePath(ctx, "exo-secrets-url")
		if err != nil {
			return nil, fmt.Errorf("resolving secrets config file path: %w", err)
		}
		secretsUrl, err := ioutil.ReadFile(secretConfigPath)
		if err == nil {
			secrets, err := ws.EsvClient.GetWorkspaceSecrets(strings.TrimSpace(string(secretsUrl)))
			if err != nil {
				ws.logEventf(ctx, "getting workspace secrets: %v", err)
			} else {
				for k, v := range secrets {
					env[k] = v
				}
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("reading secrets config: %w", err)
		}
	}

	// Encourage programs to log with colors enabled.  The closest thing to a
	// standard for this is <https://bixense.com/clicolors/#bug-reports>, but
	// support is spotty. This may grow if there are other popular enviornment
	// variables to include. If we grow PTY support, this may become unnecessary.
	env["CLICOLOR"] = "1"
	env["CLICOLOR_FORCE"] = "1"
	env["FORCE_COLOR"] = "3" // https://github.com/chalk/chalk/tree/9d5b9a133c3f8aa9f24de283660de3f732964aaa#supportscolor

	// Apply user's environment.
	// TODO: This should probably somehow shell-out to get the user's current
	// environment, otherwise changes to shell profiles won't take effect until
	// the exo daemon is exited and restarted.
	for _, assign := range os.Environ() {
		parts := strings.SplitN(assign, "=", 2)
		key := parts[0]
		val := parts[1]
		env[key] = val
	}

	// Apply .env file, if one exists.
	dotEnvMap, err := godotenv.Read(envPath)
	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("processing .env file: %w", err)
	}
	for name, value := range dotEnvMap {
		env[name] = value
	}

	return env, nil
}

func (ws *Workspace) CreateComponent(ctx context.Context, input *api.CreateComponentInput) (*api.CreateComponentOutput, error) {
	id := gensym.RandomBase32()
	jobID, err := ws.createComponent(ctx, manifest.Component{
		Name: input.Name,
		Type: input.Type,
		Spec: input.Spec,
	}, id)
	if err != nil {
		return nil, err
	}

	ws.logEventf(ctx, "new component: %s", input.Name)
	return &api.CreateComponentOutput{
		ID:    id,
		JobID: jobID,
	}, nil
}

func (ws *Workspace) createComponent(ctx context.Context, component manifest.Component, id string) (string, error) {
	if err := manifest.ValidateName(component.Name); err != nil {
		return "", errutil.HTTPErrorf(http.StatusBadRequest, "component name %q invalid: %w", component.Name, err)
	}

	if _, err := ws.Store.AddComponent(ctx, &state.AddComponentInput{
		WorkspaceID: ws.ID,
		ID:          id,
		Name:        component.Name,
		Type:        component.Type,
		Spec:        component.Spec,
		Created:     chrono.NowString(ctx),
		DependsOn:   component.DependsOn,
	}); err != nil {
		return "", fmt.Errorf("adding component: %w", err)
	}

	job := ws.TaskTracker.StartTask(contextutil.WithoutCancel(ctx), "creating "+component.Name)
	go func() {
		defer job.Finish()

		// Construct a synthetic component description to avoid re-reading after
		// the add. Only the fields needed by control are included.
		// TODO: Store.AddComponent could return a component description?
		desc := api.ComponentDescription{
			ID:        id,
			Name:      component.Name,
			Type:      component.Type,
			Spec:      component.Spec,
			DependsOn: component.DependsOn,
		}
		if err := ws.control(job, desc, &api.InitializeInput{
			Spec: component.Spec,
		}); err != nil {
			ws.logEventf(ctx, "error creating %s: %v", component.Name, err)
			job.Fail(err)
			return
		}

		if err := job.Wait(); err != nil {
			return
		}
	}()

	return job.ID(), nil
}

func (ws *Workspace) UpdateComponent(ctx context.Context, input *api.UpdateComponentInput) (*api.UpdateComponentOutput, error) {
	describeOutput, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{Refs: []string{input.Ref}})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	oldComponent := describeOutput.Components[0]
	dependsOn := oldComponent.DependsOn
	if input.DependsOn != nil {
		dependsOn = input.DependsOn
	}

	ws.logEventf(ctx, "updating %s", oldComponent.Name)
	if err := ws.updateComponent(ctx, oldComponent, manifest.Component{
		Type:      oldComponent.Type,
		Name:      oldComponent.Name,
		Spec:      input.Spec,
		DependsOn: dependsOn,
	}, oldComponent.ID); err != nil {
		return nil, err
	}

	return &api.UpdateComponentOutput{}, nil
}

func (ws *Workspace) updateComponent(ctx context.Context, oldComponent api.ComponentDescription, newComponent manifest.Component, id string) error {
	// TODO: Most updates should be accomplished without a full replacement; especially when there are no spec changes!
	if err := ws.control(ctx, oldComponent, &api.DisposeInput{}); err != nil {
		return fmt.Errorf("disposing %q for replacement: %w", oldComponent.Name, err)
	}
	if err := ws.control(ctx, oldComponent, api.InitializeInput{
		Spec: newComponent.Spec,
	}); err != nil {
		return fmt.Errorf("initializing replacement %q: %w", newComponent.Name, err)
	}
	return nil
}

func (ws *Workspace) RefreshComponents(ctx context.Context, input *api.RefreshComponentsInput) (*api.RefreshComponentsOutput, error) {
	query := makeComponentQuery(withRefs(input.Refs...))
	jobID := ws.controlEachComponent(ctx, "refreshing", query, func(desc *api.ComponentDescription) interface{} {
		return &api.RefreshInput{
			Spec: desc.Spec,
		}
	})
	return &api.RefreshComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) DisposeComponents(ctx context.Context, input *api.DisposeComponentsInput) (*api.DisposeComponentsOutput, error) {
	query := makeComponentQuery(withRefs(input.Refs...), withReversedDependencies)
	jobID := ws.controlEachComponent(ctx, "disposing", query, func(desc *api.ComponentDescription) interface{} {
		return &api.DisposeInput{}
	})
	return &api.DisposeComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) resolveRef(ctx context.Context, ref string) (string, error) {
	resolved, err := ws.resolveRefs(ctx, []string{ref})
	if err != nil {
		return "", err
	}
	return resolved[0], nil
}

func (ws *Workspace) resolveRefs(ctx context.Context, refs []string) ([]string, error) {
	resolveOutput, err := ws.Resolve(ctx, &api.ResolveInput{Refs: refs})
	if err != nil {
		return nil, err
	}
	results := make([]string, len(refs))
	for i, id := range resolveOutput.IDs {
		if id == nil {
			return nil, errutil.HTTPErrorf(http.StatusBadRequest, "unresolvable: %q", refs[i])
		}
		results[i] = *id
	}
	return results, nil
}

func (ws *Workspace) DeleteComponents(ctx context.Context, input *api.DeleteComponentsInput) (*api.DeleteComponentsOutput, error) {
	ws.logEventf(ctx, "deleting components: %s", input.Refs)
	query := makeComponentQuery(withRefs(input.Refs...), withReversedDependencies)
	jobID := ws.controlEachComponent(ctx, "deleting", query, func(*api.ComponentDescription) interface{} {
		return &api.DestroyInput{}
	})
	return &api.DeleteComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) GetComponentState(ctx context.Context, input *api.GetComponentStateInput) (*api.GetComponentStateOutput, error) {
	query := makeComponentQuery(withRefs(input.Ref))
	describe := query.describeComponentsInput(ws)

	describeOutput, err := ws.DescribeComponents(ctx, describe)
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	if len(describeOutput.Components) == 0 {
		return nil, fmt.Errorf("component not found: %q", input.Ref)
	}

	return &api.GetComponentStateOutput{
		State: describeOutput.Components[0].State,
	}, nil
}

func (ws *Workspace) SetComponentState(ctx context.Context, input *api.SetComponentStateInput) (*api.SetComponentStateOutput, error) {
	id, err := ws.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, err
	}

	// Validate that state is legal JSON.
	if !jsonutil.IsValid(input.State) {
		return nil, fmt.Errorf("state is not valid JSON")
	}

	if _, err = ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:    id,
		State: input.State,
	}); err != nil {
		return nil, fmt.Errorf("updating state: %w", err)
	}

	return &api.SetComponentStateOutput{}, nil
}

func combineLastEventAt(a, b *string) *string {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if strings.Compare(*a, *b) < 0 {
		return a
	} else {
		return b
	}
}

func (ws *Workspace) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	streamNames := input.Streams
	if streamNames == nil {
		describe := allComponentsQuery.describeComponentsInput(ws)
		components, err := ws.DescribeComponents(ctx, describe)
		if err != nil {
			return nil, fmt.Errorf("describing components: %w", err)
		}
		streamNames = make([]string, 1+len(components.Components))
		streamNames[0] = ws.ID
		for i, component := range components.Components {
			streamNames[i+1] = component.ID
		}
	}
	eventStore := log.CurrentEventStore(ctx)
	storeOutput, err := eventStore.GetEvents(ctx, &eventd.GetEventsInput{
		Streams:   streamNames,
		Cursor:    input.Cursor,
		FilterStr: input.FilterStr,
		Prev:      input.Prev,
		Next:      input.Next,
	})
	if err != nil {
		return nil, err
	}
	output := api.GetEventsOutput{
		Items:      make([]api.Event, len(storeOutput.Items)),
		PrevCursor: storeOutput.PrevCursor,
		NextCursor: storeOutput.NextCursor,
	}
	for i, storeEvent := range storeOutput.Items {
		output.Items[i] = api.Event{
			ID:        storeEvent.ID,
			Stream:    storeEvent.Stream,
			Timestamp: storeEvent.Timestamp,
			Message:   storeEvent.Message,
			Tags:      storeEvent.Tags,
		}
	}
	return &output, nil
}

func (ws *Workspace) Start(ctx context.Context, input *api.StartInput) (*api.StartOutput, error) {
	ws.logEventf(ctx, "starting...")
	jobID := ws.controlEachComponent(ctx, "starting", allProcessQuery(withDependencies), func(*api.ComponentDescription) interface{} {
		return input
	}, func(desc *api.ComponentDescription, err error) {
		ws.logEventf(ctx, "error starting %s: %v", desc.Name, err)
	})
	return &api.StartOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) StartComponents(ctx context.Context, input *api.StartComponentsInput) (*api.StartComponentsOutput, error) {
	ws.logEventf(ctx, "starting: %s", input.Refs)
	// Note that we are only querying Process component types specifically because they are the only
	// things that are "startable".
	query := allProcessQuery(withRefs(input.Refs...), withDependencies)
	jobID := ws.controlEachComponent(ctx, "starting", query, func(*api.ComponentDescription) interface{} {
		return &api.StartInput{}
	}, func(desc *api.ComponentDescription, err error) {
		ws.logEventf(ctx, "error starting %s: %v", desc.Name, err)
	})
	return &api.StartComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Stop(ctx context.Context, input *api.StopInput) (*api.StopOutput, error) {
	ws.logEventf(ctx, "stopping...")
	query := allProcessQuery(withReversedDependencies, withDependents)
	jobID := ws.controlEachComponent(ctx, "stopping", query, func(*api.ComponentDescription) interface{} {
		return input
	})
	return &api.StopOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) StopComponents(ctx context.Context, input *api.StopComponentsInput) (*api.StopComponentsOutput, error) {
	ws.logEventf(ctx, "stopping: %s", input.Refs)
	query := allProcessQuery(
		withRefs(input.Refs...),
		withReversedDependencies,
		withDependents,
	)
	jobID := ws.controlEachComponent(ctx, "stopping", query, func(*api.ComponentDescription) interface{} {
		return &api.StopInput{TimeoutSeconds: input.TimeoutSeconds}
	})
	return &api.StopComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Signal(ctx context.Context, input *api.SignalInput) (*api.SignalOutput, error) {
	ws.logEventf(ctx, "signalling %s...", input.Signal)
	query := allProcessQuery(withReversedDependencies, withDependents)
	jobID := ws.controlEachComponent(ctx, "signalling", query, func(*api.ComponentDescription) interface{} {
		return input
	})
	return &api.SignalOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) SignalComponents(ctx context.Context, input *api.SignalComponentsInput) (*api.SignalComponentsOutput, error) {
	ws.logEventf(ctx, "signalling %s to ", input.Refs)
	query := allProcessQuery(
		withRefs(input.Refs...),
		withReversedDependencies,
		withDependents,
	)
	jobID := ws.controlEachComponent(ctx, "signalling", query, func(*api.ComponentDescription) interface{} {
		return &api.SignalInput{
			Signal: input.Signal,
		}
	})
	return &api.SignalComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Restart(ctx context.Context, input *api.RestartInput) (*api.RestartOutput, error) {
	ws.logEventf(ctx, "restarting...")
	query := makeComponentQuery(withDependencies)
	jobID := ws.controlEachComponent(ctx, "restarting", query, func(*api.ComponentDescription) interface{} {
		return input
	}, func(desc *api.ComponentDescription, err error) {
		ws.logEventf(ctx, "error restarting %s: %v", desc.Name, err)
	})
	return &api.RestartOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) RestartComponents(ctx context.Context, input *api.RestartComponentsInput) (*api.RestartComponentsOutput, error) {
	ws.logEventf(ctx, "restarting %s", input.Refs)
	// Restart currently restarts the component and everything that depends on it in the same order as start. There
	// are likely 3 different restart "modes" that we will eventually want to support:
	// 1. Restart only the component(s) requested. Do not cascade restarts to any other components.
	// 2. Stop everything that depends on these components in reverse dependency order, restart these components,
	//    then restart the dependents in normal order again.
	// 3. Ensure that all dependendencies are started (in normal order), then restart these components (current behaviour).
	query := allProcessQuery(withRefs(input.Refs...), withDependencies)
	jobID := ws.controlEachComponent(ctx, "restart", query, func(*api.ComponentDescription) interface{} {
		return &api.RestartInput{TimeoutSeconds: input.TimeoutSeconds}
	}, func(desc *api.ComponentDescription, err error) {
		ws.logEventf(ctx, "error restarting %s: %v", desc.Name, err)
	})
	return &api.RestartComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) DescribeProcesses(ctx context.Context, input *api.DescribeProcessesInput) (*api.DescribeProcessesOutput, error) {
	describe := allProcessQuery().describeComponentsInput(ws)
	components, err := ws.DescribeComponents(ctx, describe)
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	var eg errgroup.Group
	processes := make([]api.ProcessDescription, len(components.Components))
	for i, component := range components.Components {
		i, component := i, component
		eg.Go(func() error {
			var desc api.ProcessDescription
			var err error
			// XXX Violates component state encapsulation.
			switch component.Type {
			case "process":
				desc, err = process.GetProcessDescription(ctx, component)
			case "container":
				desc, err = container.GetProcessDescription(ctx, ws.Docker, component)
			}
			if err != nil {
				return fmt.Errorf("could not get process description: %w", err)
			}
			processes[i] = desc
			return nil
		})
	}
	err = eg.Wait()
	return &api.DescribeProcessesOutput{Processes: processes}, err
}

func (ws *Workspace) DescribeVolumes(ctx context.Context, input *api.DescribeVolumesInput) (*api.DescribeVolumesOutput, error) {
	query := makeComponentQuery(withTypes("volume"))
	describe := query.describeComponentsInput(ws)
	components, err := ws.DescribeComponents(ctx, describe)
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	output := api.DescribeVolumesOutput{
		Volumes: make([]api.VolumeDescription, 0, len(components.Components)),
	}
	for _, component := range components.Components {
		volume := api.VolumeDescription{
			ID:   component.ID,
			Name: component.Name,
		}
		output.Volumes = append(output.Volumes, volume)
	}
	return &output, nil
}

func (ws *Workspace) DescribeNetworks(ctx context.Context, input *api.DescribeNetworksInput) (*api.DescribeNetworksOutput, error) {
	query := makeComponentQuery(withTypes("network"))
	describe := query.describeComponentsInput(ws)
	components, err := ws.DescribeComponents(ctx, describe)
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	output := api.DescribeNetworksOutput{
		Networks: make([]api.NetworkDescription, 0, len(components.Components)),
	}
	for _, component := range components.Components {
		network := api.NetworkDescription{
			ID:   component.ID,
			Name: component.Name,
		}
		output.Networks = append(output.Networks, network)
	}
	return &output, nil
}

func (ws *Workspace) ExportProcfile(ctx context.Context, input *api.ExportProcfileInput) (*api.ExportProcfileOutput, error) {
	logger := logging.CurrentLogger(ctx)
	procs, err := ws.DescribeProcesses(ctx, &api.DescribeProcessesInput{})
	if err != nil {
		return nil, fmt.Errorf("describing processes: %w", err)
	}

	unixProcs := make([]procfile.Process, 0, len(procs.Processes))
	for _, proc := range procs.Processes {
		if proc.Provider == "unix" {
			var spec process.Spec
			if err := jsonutil.UnmarshalStringOrEmpty(proc.Spec, &spec); err != nil {
				logger.Infof("unmarshalling process spec: %v\n", err)
				continue
			}

			unixProcs = append(unixProcs, procfile.Process{
				Name:        proc.Name,
				Program:     spec.Program,
				Arguments:   spec.Arguments,
				Environment: spec.Environment,
			})
		}
	}

	// Produce a stable order of processes for export.  Ideally, this would
	// preserve the original order from an imported Procfile, but that would
	// require some metadata on components.  Acheiving a stable order on the
	// first export is the next best thing.
	procfile.Organize(&unixProcs)

	var export bytes.Buffer
	if err := procfile.Generate(&export, unixProcs); err != nil {
		return nil, fmt.Errorf("generating procfile: %w", err)
	}

	return &api.ExportProcfileOutput{
		Procfile: export.String(),
	}, nil
}

func (ws *Workspace) ReadFile(ctx context.Context, input *api.ReadFileInput) (*api.ReadFileOutput, error) {
	resolvedPath, err := ws.resolveWorkspacePath(ctx, input.Path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	return &api.ReadFileOutput{
		Content: string(content),
	}, nil
}

func (ws *Workspace) WriteFile(ctx context.Context, input *api.WriteFileInput) (*api.WriteFileOutput, error) {
	resolvedPath, err := ws.resolveWorkspacePath(ctx, input.Path)
	if err != nil {
		return nil, err
	}

	mode := os.FileMode(0644)
	if input.Mode != nil {
		mode = os.FileMode(*input.Mode)
	}

	if err := os.WriteFile(resolvedPath, []byte(input.Content), mode); err != nil {
		return nil, err
	}

	return &api.WriteFileOutput{}, nil
}

func (ws *Workspace) resolveWorkspacePath(ctx context.Context, relativePath string) (string, error) {
	description, err := ws.describe(ctx)
	if err != nil {
		return "", fmt.Errorf("describing workspace: %w", err)
	}

	resolvedPath := path.Join(description.Root, relativePath)
	if !pathutil.HasPathPrefix(resolvedPath, description.Root) {
		return "", fmt.Errorf("directory is outside workspace: %q", relativePath)
	}

	return resolvedPath, nil
}

func (ws *Workspace) controlEachComponent(ctx context.Context, label string, query componentQuery, makeMessage func(*api.ComponentDescription) interface{}, onErr ...func(*api.ComponentDescription, error)) (jobID string) {
	job := ws.TaskTracker.StartTask(ctx, label)
	go func() {
		defer job.Finish()
		ws.goControlComponents(job, query, makeMessage, onErr...)
	}()
	return job.ID()
}

func (ws *Workspace) goControlComponents(t *task.Task, query componentQuery, makeMessage func(*api.ComponentDescription) interface{}, onErr ...func(*api.ComponentDescription, error)) {
	describe := query.describeComponentsInput(ws)
	components, err := ws.DescribeComponents(t, describe)
	if err != nil {
		t.Fail(fmt.Errorf("describing components: %w", err))
		return
	}

	// Build graph of tasks to run.
	runGraph := deps.New()
	for _, component := range components.Components {
		component := component
		runGraph.AddNode(&runTaskNode{
			name: component.Name,
			task: t.CreateChild(component.Name),
			run: func(t *task.Task) error {
				err := ws.control(t, component, makeMessage(&component))
				if err != nil {
					for _, f := range onErr {
						f(&component, err)
					}
				}
				return err
			},
		})
		for _, dependency := range component.DependsOn {
			runGraph.AddEdge(component.Name, dependency)
		}
	}

	// Run tasks.
	layers := runGraph.TopoSortedLayers()
	for i, layer := range layers {
		if query.DependencyOrder == dependencyOrderReverse {
			layer = layers[len(layers)-1-i]
		}

		var wg sync.WaitGroup
		for _, node := range layer {
			runTask := node.(*runTaskNode)
			wg.Add(1)
			go func() {
				defer wg.Done()
				runTask.task.Start()
				defer runTask.task.Finish()
				if err := runTask.run(runTask.task); err != nil {
					runTask.task.Fail(err)
				}
			}()
		}
		wg.Wait()
	}
}

func (ws *Workspace) control(ctx context.Context, desc api.ComponentDescription, input interface{}) error {
	ctrl := ws.newController(ctx, desc)
	if err := ctrl.InitResource(); err != nil {
		return err
	}
	// TODO: Figure out how to avoid special-casing destroy.
	destroying := false
	switch input.(type) {
	case *api.DestroyInput:
		destroying = true
		input = &api.DisposeInput{}
	}
	_, fErr := josh.Send(ctx, ctrl, input)
	// Try to save state even if f fails.
	newState, err := ctrl.MarshalState()
	if err == nil {
		if destroying {
			// TODO: Do this asynchronously as a garbage collection pass, so
			// that we can inspect deleted components and debug them.
			_, err = ws.Store.RemoveComponent(ctx, &state.RemoveComponentInput{
				ID: desc.ID,
			})
		} else {
			_, err = ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
				ID:    desc.ID,
				State: newState,
			})
		}
	}
	if fErr != nil {
		return fErr
	}
	return err
}

func (ws *Workspace) Build(ctx context.Context, input *api.BuildInput) (*api.BuildOutput, error) {
	ws.logEventf(ctx, "building...")
	jobID := ws.controlEachComponent(ctx, "building", allBuildableQuery(), func(*api.ComponentDescription) interface{} {
		return input
	})
	return &api.BuildOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) BuildComponents(ctx context.Context, input *api.BuildComponentsInput) (*api.BuildComponentsOutput, error) {
	ws.logEventf(ctx, "building: %s", input.Refs)
	query := allBuildableQuery(withRefs(input.Refs...))
	jobID := ws.controlEachComponent(ctx, "building", query, func(*api.ComponentDescription) interface{} {
		return &api.BuildInput{}
	})
	return &api.BuildComponentsOutput{
		JobID: jobID,
	}, nil
}

type runTaskNode struct {
	name string
	task *task.Task
	run  func(task *task.Task) error
}

func (n *runTaskNode) ID() string {
	return n.name
}

type componentNode struct {
	component manifest.Component
}

func (n *componentNode) ID() string {
	return n.component.Name
}
