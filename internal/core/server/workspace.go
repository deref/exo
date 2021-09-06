package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"

	"github.com/deref/exo/internal/chrono"
	"github.com/deref/exo/internal/core/api"
	state "github.com/deref/exo/internal/core/state/api"
	"github.com/deref/exo/internal/deps"
	"github.com/deref/exo/internal/gensym"
	logd "github.com/deref/exo/internal/logd/api"
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
	return &api.WorkspaceDescription{
		ID:   ws.ID,
		Root: output.Workspaces[0].Root,
	}, nil
}

func (ws *Workspace) Destroy(ctx context.Context, input *api.DestroyInput) (*api.DestroyOutput, error) {
	job := ws.TaskTracker.StartTask(ctx, "destroying")
	query := makeComponentQuery(withReversedDependencies)

	go func() {
		defer job.Finish()
		ws.goControlComponents(job, query, func(ctx context.Context, lifecycle api.Lifecycle) error {
			return ws.deleteComponent(ctx, lifecycle)
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
				return ws.control(job.Context, oldComponent, func(ctx context.Context, lifecycle api.Lifecycle) error {
					return ws.deleteComponent(job.Context, lifecycle)
				})
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
				return ws.control(job.Context, oldComponent, func(ctx context.Context, lifecycle api.Lifecycle) error {
					return ws.deleteComponent(job.Context, lifecycle)
				})
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
				return ws.createComponent(t, newComponent, gensym.RandomBase32())
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
					return ws.createComponent(t, newComponent, id)
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
		IDs:                 input.IDs,
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
			ID:          component.ID,
			Name:        component.Name,
			Type:        component.Type,
			Spec:        component.Spec,
			State:       component.State,
			Created:     component.Created,
			Initialized: component.Initialized,
			Disposed:    component.Disposed,
			DependsOn:   component.DependsOn,
		}
	}
	return output, nil
}

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
		ComponentSpec:        desc.Spec,
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
		return nil, fmt.Errorf("could not resolve path to .env file: %w", err)
	}

	_, err = os.Stat(envPath)
	if os.IsNotExist(err) {
		env := make(map[string]string)
		for _, assign := range os.Environ() {
			parts := strings.SplitN(assign, "=", 2)
			key := parts[0]
			val := parts[1]
			env[key] = val
		}
		return env, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not read .env file: %w", err)
	}

	envMap, err := godotenv.Read(envPath)
	if err != nil {
		return nil, fmt.Errorf("could not process env file: %w", err)
	}
	return envMap, nil
}

func (ws *Workspace) CreateComponent(ctx context.Context, input *api.CreateComponentInput) (*api.CreateComponentOutput, error) {
	id := gensym.RandomBase32()
	err := ws.createComponent(ctx, manifest.Component{
		Name: input.Name,
		Type: input.Type,
		Spec: input.Spec,
	}, id)
	if err != nil {
		return nil, err
	}
	return &api.CreateComponentOutput{
		ID: id,
	}, nil
}

func (ws *Workspace) createComponent(ctx context.Context, component manifest.Component, id string) error {
	if err := manifest.ValidateName(component.Name); err != nil {
		return errutil.HTTPErrorf(http.StatusBadRequest, "component name %q invalid: %w", component.Name, err)
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
		return fmt.Errorf("adding component: %w", err)
	}

	if err := ws.control(ctx, api.ComponentDescription{
		// Construct a synthetic component description to avoid re-reading after
		// the add. Only the fields needed by control are included.
		// TODO: Store.AddComponent could return a component description?
		ID:        id,
		Name:      component.Name,
		Type:      component.Type,
		Spec:      component.Spec,
		DependsOn: component.DependsOn,
	}, func(ctx context.Context, lifecycle api.Lifecycle) error {
		_, err := lifecycle.Initialize(ctx, &api.InitializeInput{})
		return err
	}); err != nil {
		return err
	}

	// XXX this now double-patches the component to set Initialized timestamp. Optimize?
	if _, err := ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:          id,
		Initialized: chrono.NowString(ctx),
	}); err != nil {
		return fmt.Errorf("modifying component after initialization: %w", err) // XXX this message seems incorrect.
	}

	return nil
}

func (ws *Workspace) UpdateComponent(ctx context.Context, input *api.UpdateComponentInput) (*api.UpdateComponentOutput, error) {
	id, err := ws.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, err
	}

	describeOutput, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{IDs: []string{id}})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	oldComponent := describeOutput.Components[0]
	dependsOn := oldComponent.DependsOn
	if input.DependsOn != nil {
		dependsOn = input.DependsOn
	}

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
	if err := ws.control(ctx, oldComponent, func(ctx context.Context, lifecycle api.Lifecycle) error {
		return ws.deleteComponent(ctx, lifecycle)
	}); err != nil {
		return fmt.Errorf("delete %q for replacement: %w", oldComponent.Name, err)
	}
	if err := ws.createComponent(ctx, newComponent, id); err != nil {
		return fmt.Errorf("adding replacement %q: %w", newComponent.Name, err)
	}
	return nil
}

func (ws *Workspace) RefreshComponents(ctx context.Context, input *api.RefreshComponentsInput) (*api.RefreshComponentsOutput, error) {
	query := makeComponentQuery(withRefs(input.Refs...))
	jobID := ws.controlEachComponent(ctx, "refreshing", query, func(ctx context.Context, lifecycle api.Lifecycle) error {
		_, err := lifecycle.Refresh(ctx, &api.RefreshInput{})
		return err
	})
	return &api.RefreshComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) DisposeComponents(ctx context.Context, input *api.DisposeComponentsInput) (*api.DisposeComponentsOutput, error) {
	query := makeComponentQuery(withRefs(input.Refs...), withReversedDependencies)
	jobID := ws.controlEachComponent(ctx, "disposing", query, func(ctx context.Context, lifecycle api.Lifecycle) error {
		_, err := lifecycle.Dispose(ctx, &api.DisposeInput{})
		return err
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
	query := makeComponentQuery(withRefs(input.Refs...), withReversedDependencies)
	jobID := ws.controlEachComponent(ctx, "deleting", query, func(ctx context.Context, lifecycle api.Lifecycle) error {
		return ws.deleteComponent(ctx, lifecycle)
	})
	return &api.DeleteComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) deleteComponent(ctx context.Context, lifecycle api.Lifecycle) error {
	_, err := lifecycle.Dispose(ctx, &api.DisposeInput{})
	if err != nil {
		return err
	}
	lifecycle.(core.Component).MarkDeleted()
	return nil
}

func (ws *Workspace) GetComponentState(ctx context.Context, input *api.GetComponentStateInput) (*api.GetComponentStateOutput, error) {
	query := makeComponentQuery(withRefs(input.Ref))
	describe, err := query.describeComponentsInput(ctx, ws)
	if err != nil {
		return nil, err
	}

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

func (ws *Workspace) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
	describe, err := allProcessQuery().describeComponentsInput(ctx, ws)
	if err != nil {
		return nil, err
	}

	components, err := ws.DescribeComponents(ctx, describe)
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	// Find all logs in component hierarchy.
	// TODO: More general handling of log groups, subcomponents, etc.
	var logGroups []string
	var logStreams []string
	streamToGroup := make(map[string]int)
	for _, component := range components.Components {
		// XXX Janky provider inference. See note: [LOG_COMPONENTS].
		var provider string
		switch component.Type {
		case "process":
			provider = "unix"
		case "container":
			provider = "docker"
		}
		for _, streamName := range log.ComponentLogNames(provider, component.ID) {
			streamToGroup[streamName] = len(logGroups)
			logStreams = append(logStreams, streamName)
		}
		logGroups = append(logGroups, component.ID)
	}

	// Initialize output and index by log group name.
	logs := make([]api.LogDescription, len(logGroups))
	for i, logGroup := range logGroups {
		logs[i] = api.LogDescription{
			Name: logGroup,
		}
	}

	// Decorate output with information from the log collector.
	collector := log.CurrentLogCollector(ctx)
	collectorLogs, err := collector.DescribeLogs(ctx, &logd.DescribeLogsInput{
		Names: logStreams,
	})
	if err != nil {
		return nil, err
	}
	for _, collectorLog := range collectorLogs.Logs {
		groupIndex, ok := streamToGroup[collectorLog.Name]
		if !ok {
			continue
		}
		group := &logs[groupIndex]
		group.LastEventAt = combineLastEventAt(group.LastEventAt, collectorLog.LastEventAt)
	}
	return &api.DescribeLogsOutput{Logs: logs}, nil
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
	logGroups := input.Logs
	if logGroups == nil {
		// No filter specified, use all streams.
		logDescriptions, err := ws.DescribeLogs(ctx, &api.DescribeLogsInput{})
		if err != nil {
			return nil, fmt.Errorf("enumerating logs: %w", err)
		}
		logGroups = make([]string, len(logDescriptions.Logs))
		for i, group := range logDescriptions.Logs {
			logGroups[i] = group.Name
		}
	}
	logStreams := make([]string, 0, 2*len(logGroups))
	// Expand log groups in to streams.
	for _, group := range logGroups {
		// Each process acts as a log group combining both stdout and stderr.
		// XXX See note [LOG_COMPONENTS].
		for _, suffix := range []string{"", ":out", ":err"} {
			stream := group + suffix
			logStreams = append(logStreams, stream)
		}
	}

	collector := log.CurrentLogCollector(ctx)
	collectorOutput, err := collector.GetEvents(ctx, &logd.GetEventsInput{
		Logs:      logStreams,
		Cursor:    input.Cursor,
		FilterStr: input.FilterStr,
		Prev:      input.Prev,
		Next:      input.Next,
	})
	if err != nil {
		return nil, err
	}
	output := api.GetEventsOutput{
		Items:      make([]api.Event, len(collectorOutput.Items)),
		PrevCursor: collectorOutput.PrevCursor,
		NextCursor: collectorOutput.NextCursor,
	}
	for i, collectorEvent := range collectorOutput.Items {
		output.Items[i] = api.Event{
			ID:        collectorEvent.ID,
			Log:       collectorEvent.Log,
			Timestamp: collectorEvent.Timestamp,
			Message:   collectorEvent.Message,
		}
	}
	return &output, nil
}

func (ws *Workspace) Start(ctx context.Context, input *api.StartInput) (*api.StartOutput, error) {
	jobID := ws.controlEachComponent(ctx, "starting", allProcessQuery(withDependencies), func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Start(ctx, &api.StartInput{})
			return err
		}
		if lifecycle, ok := thing.(api.Lifecycle); ok {
			_, err := lifecycle.Initialize(ctx, &api.InitializeInput{})
			return err
		}
		return nil
	})
	return &api.StartOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) StartComponents(ctx context.Context, input *api.StartComponentsInput) (*api.StartComponentsOutput, error) {
	// Note that we are only querying Process component types specifically because they are the only
	// things that are "startable".
	query := allProcessQuery(withRefs(input.Refs...), withDependencies)
	jobID := ws.controlEachComponent(ctx, "starting", query, func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Start(ctx, &api.StartInput{})
			return err
		}
		if lifecycle, ok := thing.(api.Lifecycle); ok {
			_, err := lifecycle.Initialize(ctx, &api.InitializeInput{})
			return err
		}
		return nil
	})
	return &api.StartComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Stop(ctx context.Context, input *api.StopInput) (*api.StopOutput, error) {
	query := allProcessQuery(withReversedDependencies, withDependents)
	jobID := ws.controlEachComponent(ctx, "stopping", query, func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Stop(ctx, input)
			return err
		}
		return nil
	})
	return &api.StopOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) StopComponents(ctx context.Context, input *api.StopComponentsInput) (*api.StopComponentsOutput, error) {
	query := allProcessQuery(
		withRefs(input.Refs...),
		withReversedDependencies,
		withDependents,
	)
	jobID := ws.controlEachComponent(ctx, "stopping", query, func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Stop(ctx, &api.StopInput{TimeoutSeconds: input.TimeoutSeconds})
			return err
		}
		return nil
	})
	return &api.StopComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Signal(ctx context.Context, input *api.SignalInput) (*api.SignalOutput, error) {
	query := allProcessQuery(withReversedDependencies, withDependents)
	jobID := ws.controlEachComponent(ctx, "signalling", query, func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Signal(ctx, input)
			return err
		}
		return nil
	})
	return &api.SignalOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) SignalComponents(ctx context.Context, input *api.SignalComponentsInput) (*api.SignalComponentsOutput, error) {
	query := allProcessQuery(
		withRefs(input.Refs...),
		withReversedDependencies,
		withDependents,
	)
	jobID := ws.controlEachComponent(ctx, "signalling", query, func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Signal(ctx, &api.SignalInput{
				Signal: input.Signal,
			})
			return err
		}
		return nil
	})
	return &api.SignalComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) Restart(ctx context.Context, input *api.RestartInput) (*api.RestartOutput, error) {
	query := makeComponentQuery(withDependencies)
	jobID := ws.controlEachComponent(ctx, "restarting", query, func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Restart(ctx, input)
			return err
		}
		return nil
	})
	return &api.RestartOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) RestartComponents(ctx context.Context, input *api.RestartComponentsInput) (*api.RestartComponentsOutput, error) {
	// Restart currently restarts the component and everything that depends on it in the same order as start. There
	// are likely 3 different restart "modes" that we will eventually want to support:
	// 1. Restart only the component(s) requested. Do not cascade restarts to any other components.
	// 2. Stop everything that depends on these components in reverse dependency order, restart these components,
	//    then restart the dependents in normal order again.
	// 3. Ensure that all dependendencies are started (in normal order), then restart these components (current behaviour).
	query := allProcessQuery(withRefs(input.Refs...), withDependencies)
	jobID := ws.controlEachComponent(ctx, "restart", query, func(ctx context.Context, thing interface{}) error {
		if process, ok := thing.(api.Process); ok {
			_, err := process.Restart(ctx, &api.RestartInput{TimeoutSeconds: input.TimeoutSeconds})
			return err
		}
		return nil
	})
	return &api.RestartComponentsOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) DescribeProcesses(ctx context.Context, input *api.DescribeProcessesInput) (*api.DescribeProcessesOutput, error) {
	describe, err := allProcessQuery().describeComponentsInput(ctx, ws)
	if err != nil {
		return nil, err
	}

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
	describe, err := query.describeComponentsInput(ctx, ws)
	if err != nil {
		return nil, err
	}
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
	describe, err := query.describeComponentsInput(ctx, ws)
	if err != nil {
		return nil, err
	}
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
			if err := jsonutil.UnmarshalString(proc.Spec, &spec); err != nil {
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

func (ws *Workspace) controlEachComponent(ctx context.Context, label string, query componentQuery, f interface{}) (jobID string) {
	job := ws.TaskTracker.StartTask(ctx, label)
	go func() {
		defer job.Finish()
		ws.goControlComponents(job, query, f)
	}()
	return job.ID()
}

func (ws *Workspace) goControlComponents(t *task.Task, query componentQuery, f interface{}) {
	describe, err := query.describeComponentsInput(t, ws)
	if err != nil {
		t.Fail(err)
		return
	}

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
				return ws.control(t, component, f)
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

func (ws *Workspace) control(ctx context.Context, desc api.ComponentDescription, f interface{}) error {
	ctrl := ws.newController(ctx, desc)
	if err := ctrl.InitResource(); err != nil {
		return err
	}
	ctxV := reflect.ValueOf(ctx)
	fV := reflect.ValueOf(f)
	ctrlV := reflect.ValueOf(ctrl)
	argT := fV.Type().In(1)
	if !ctrlV.Type().AssignableTo(argT) {
		return fmt.Errorf("%q controller does not implement %s", desc.Type, argT)
	}
	results := fV.Call([]reflect.Value{ctxV, ctrlV})
	fErr, _ := results[0].Interface().(error)
	// Try to save state even if f fails.
	newState, err := ctrl.MarshalState()
	if err == nil {
		if ctrl.IsDeleted() {
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
	jobID := ws.controlEachComponent(ctx, "building", allBuildableQuery(), func(ctx context.Context, builder api.Builder) error {
		_, err := builder.Build(ctx, &api.BuildInput{})
		return err
	})
	return &api.BuildOutput{
		JobID: jobID,
	}, nil
}

func (ws *Workspace) BuildComponents(ctx context.Context, input *api.BuildComponentsInput) (*api.BuildComponentsOutput, error) {
	query := allBuildableQuery(withRefs(input.Refs...))
	jobID := ws.controlEachComponent(ctx, "building", query, func(ctx context.Context, builder api.Builder) error {
		_, err := builder.Build(ctx, &api.BuildInput{})
		return err
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
