package server

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/components/invalid"
	"github.com/deref/exo/components/log"
	"github.com/deref/exo/components/process"
	"github.com/deref/exo/config"
	core "github.com/deref/exo/core/api"
	"github.com/deref/exo/exod/api"
	state "github.com/deref/exo/exod/state/api"
	"github.com/deref/exo/gensym"
	"github.com/deref/exo/import/procfile"
	logd "github.com/deref/exo/logd/api"
	"github.com/deref/exo/util/jsonutil"
)

type Workspace struct {
	ID     string
	VarDir string
	Store  state.Store
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
	describeOutput, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	// TODO: Parallelism / bulk delete.
	for _, component := range describeOutput.Components {
		if err := ws.deleteComponent(ctx, component.ID); err != nil {
			return nil, fmt.Errorf("deleting %s: %w", component.Name, err)
		}
	}
	if _, err := ws.Store.RemoveWorkspace(ctx, &state.RemoveWorkspaceInput{
		ID: ws.ID,
	}); err != nil {
		return nil, fmt.Errorf("removing workspace from store: %w", err)
	}
	return &api.DestroyOutput{}, nil
}

func (ws *Workspace) Apply(ctx context.Context, input *api.ApplyInput) (*api.ApplyOutput, error) {
	cfg, err := config.Parse([]byte(input.Config))
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if err := ws.apply(ctx, cfg); err != nil {
		return nil, err
	}
	return &api.ApplyOutput{}, nil
}

func (ws *Workspace) apply(ctx context.Context, cfg *config.Config) error {
	// TODO: Validate config.

	describeOutput, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
	})
	if err != nil {
		return fmt.Errorf("describing components: %w", err)
	}

	// Index old components by name.
	oldComponents := make(map[string]state.ComponentDescription, len(describeOutput.Components))
	for _, component := range describeOutput.Components {
		oldComponents[component.Name] = component
	}

	// TODO: Handle partial failures.

	// Apply component upserts.
	newComponents := make(map[string]config.Component, len(cfg.Components))
	for _, newComponent := range cfg.Components {
		name := newComponent.Name
		newComponents[name] = newComponent
		if oldComponent, exists := oldComponents[name]; exists {
			// Update existing component.
			if err := ws.updateComponent(ctx, oldComponent, newComponent); err != nil {
				return fmt.Errorf("updating %q: %w", name, err)
			}
		} else {
			// Create new component.
			if _, err := ws.createComponent(ctx, newComponent); err != nil {
				return fmt.Errorf("adding %q: %w", name, err)
			}
		}
	}

	// Apply component deletions.
	// TODO: Dispose in parallel. Optionally await deletion.
	for name, oldComponent := range oldComponents {
		if _, keep := newComponents[name]; keep {
			continue
		}
		if err := ws.deleteComponent(ctx, oldComponent.ID); err != nil {
			return fmt.Errorf("deleting %q: %w", name, err)
		}
	}

	return nil
}

func (ws *Workspace) ApplyProcfile(ctx context.Context, input *api.ApplyProcfileInput) (*api.ApplyProcfileOutput, error) {
	cfg, err := procfile.Import(strings.NewReader(input.Procfile))
	if err != nil {
		return nil, fmt.Errorf("importing: %w", err)
	}
	if err := ws.apply(ctx, cfg); err != nil {
		return nil, err
	}
	return &api.ApplyProcfileOutput{}, nil
}

func (ws *Workspace) Refresh(ctx context.Context, input *api.RefreshInput) (*api.RefreshOutput, error) {
	components, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
	})
	if err != nil {
		return nil, err
	}
	// TODO: Parallelism.
	for _, component := range components.Components {
		if err := ws.refreshComponent(ctx, ws.Store, component); err != nil {
			// TODO: Error recovery.
			return nil, fmt.Errorf("refreshing %q: %w", component.Name, err)
		}
	}
	return &api.RefreshOutput{}, nil
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
		WorkspaceID: ws.ID,
	})
	if err != nil {
		return nil, err
	}
	output := &api.DescribeComponentsOutput{
		Components: []api.ComponentDescription{},
	}
	for _, component := range stateOutput.Components {
		output.Components = append(output.Components, api.ComponentDescription{
			ID:          component.ID,
			Name:        component.Name,
			Type:        component.Type,
			Spec:        component.Spec,
			State:       component.State,
			Created:     component.Created,
			Initialized: component.Initialized,
			Disposed:    component.Disposed,
		})
	}
	return output, nil
}

func (ws *Workspace) resolveProvider(ctx context.Context, typ string) core.Provider {
	switch typ {
	case "process":
		description, err := ws.describe(ctx)
		if err != nil {
			return &invalid.Provider{
				fmt.Errorf("workspace error: %w", err),
			}
		}
		return &process.Provider{
			WorkspaceDir: description.Root,
			VarDir:       filepath.Join(ws.VarDir, "proc"),
		}
	default:
		return &invalid.Provider{
			Err: fmt.Errorf("unsupported component type: %q", typ),
		}
	}
}

func (ws *Workspace) CreateComponent(ctx context.Context, input *api.CreateComponentInput) (*api.CreateComponentOutput, error) {
	id, err := ws.createComponent(ctx, config.Component{
		Name: input.Name,
		Type: input.Type,
		Spec: input.Spec,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreateComponentOutput{
		ID: id,
	}, nil
}

func (ws *Workspace) createComponent(ctx context.Context, component config.Component) (id string, err error) {
	if !IsValidName(component.Name) {
		return "", fmt.Errorf("invalid name: %q", component.Name)
	}

	id = gensym.RandomBase32()

	if _, err := ws.Store.AddComponent(ctx, &state.AddComponentInput{
		WorkspaceID: ws.ID,
		ID:          id,
		Name:        component.Name,
		Type:        component.Type,
		Spec:        component.Spec,
		Created:     chrono.NowString(ctx),
	}); err != nil {
		return "", fmt.Errorf("adding component: %w", err)
	}

	provider := ws.resolveProvider(ctx, component.Type)
	output, err := provider.Initialize(ctx, &core.InitializeInput{
		ID:   id,
		Spec: component.Spec,
	})
	if err != nil {
		return "", err
	}

	if _, err := ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:          id,
		State:       output.State,
		Initialized: chrono.NowString(ctx),
	}); err != nil {
		return "", fmt.Errorf("modifying component after initialization: %w", err)
	}

	return id, nil
}

func IsValidName(name string) bool {
	for _, b := range []byte(name) {
		if b == 0 || b == 255 {
			return false
		}
	}
	return name != ""
}

func (ws *Workspace) UpdateComponent(ctx context.Context, input *api.UpdateComponentInput) (*api.UpdateComponentOutput, error) {
	panic("TODO: UpdateComponent")
}

func (ws *Workspace) updateComponent(ctx context.Context, oldComponent state.ComponentDescription, newComponent config.Component) error {
	// TODO: Smart updating, using update lifecycle method.
	name := oldComponent.Name
	id := oldComponent.ID
	if err := ws.deleteComponent(ctx, id); err != nil {
		return fmt.Errorf("delete %q for replacement: %w", name, err)
	}
	if _, err := ws.createComponent(ctx, newComponent); err != nil {
		return fmt.Errorf("adding replacement %q: %w", name, err)
	}
	return nil
}

func (ws *Workspace) RefreshComponent(ctx context.Context, input *api.RefreshComponentInput) (*api.RefreshComponentOutput, error) {
	id, err := ws.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving ref: %w", err)
	}

	describeOutput, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
		IDs:         []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	if len(describeOutput.Components) < 1 {
		return nil, fmt.Errorf("no component %q", id)
	}
	component := describeOutput.Components[0]

	if err := ws.refreshComponent(ctx, ws.Store, component); err != nil {
		return nil, err
	}

	return &api.RefreshComponentOutput{}, err
}

func (ws *Workspace) refreshComponent(ctx context.Context, store state.Store, component state.ComponentDescription) error {
	provider := ws.resolveProvider(ctx, component.Type)
	refreshed, err := provider.Refresh(ctx, &core.RefreshInput{
		ID:    component.ID,
		Spec:  component.Spec,
		State: component.State,
	})
	if err != nil {
		return err
	}

	if _, err := store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:          component.ID,
		State:       refreshed.State,
		Initialized: chrono.NowString(ctx),
	}); err != nil {
		return fmt.Errorf("modifying component after initialization: %w", err)
	}
	return nil
}

func (ws *Workspace) DisposeComponent(ctx context.Context, input *api.DisposeComponentInput) (*api.DisposeComponentOutput, error) {
	id, err := ws.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving ref: %w", err)
	}
	err = ws.disposeComponent(ctx, id)
	return &api.DisposeComponentOutput{}, err
}

func (ws *Workspace) resolveRef(ctx context.Context, ref string) (string, error) {
	resolveOutput, err := ws.Resolve(ctx, &api.ResolveInput{Refs: []string{ref}})
	if err != nil {
		return "", err
	}
	id := resolveOutput.IDs[0]
	if id == nil {
		return "", errors.New("unresolvable")
	}
	return *id, nil
}

func (ws *Workspace) disposeComponent(ctx context.Context, id string) error {
	describeOutput, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
		IDs:         []string{id},
	})
	if err != nil {
		return fmt.Errorf("describing components: %w", err)
	}
	if len(describeOutput.Components) < 1 {
		return fmt.Errorf("no component %q", id)
	}
	component := describeOutput.Components[0]
	provider := ws.resolveProvider(ctx, component.Type)
	_, err = provider.Dispose(ctx, &core.DisposeInput{
		ID:    id,
		State: component.State,
	})
	return err
}

func (ws *Workspace) DeleteComponent(ctx context.Context, input *api.DeleteComponentInput) (*api.DeleteComponentOutput, error) {
	id, err := ws.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving ref: %w", err)
	}
	if err := ws.deleteComponent(ctx, id); err != nil {
		return nil, err
	}
	return &api.DeleteComponentOutput{}, nil
}

func (ws *Workspace) deleteComponent(ctx context.Context, id string) error {
	if err := ws.disposeComponent(ctx, id); err != nil {
		return fmt.Errorf("disposing: %w", err)
	}
	// TODO: Await disposal.
	if _, err := ws.Store.RemoveComponent(ctx, &state.RemoveComponentInput{ID: id}); err != nil {
		return fmt.Errorf("removing from state store: %w", err)
	}
	return nil
}

var processLogStreams = []string{"out", "err"}

func (ws *Workspace) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
	components, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}

	// Find all logs in component hierarchy.
	// TODO: More general handling of log groups, subcomponents, etc.
	var logGroups []string
	var logStreams []string
	streamToGroup := make(map[string]int, len(logGroups))
	for _, component := range components.Components {
		if component.Type == "process" {
			for _, stream := range processLogStreams {
				streamName := fmt.Sprintf("%s:%s", component.ID, stream)
				streamToGroup[streamName] = len(logGroups)
				logStreams = append(logStreams, streamName)
			}
			logGroups = append(logGroups, component.ID)
		}
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
	var logStreams []string
	// Expand log groups in to streams.
	for _, group := range logGroups {
		// Each process acts as a log group combining both stdout and stderr.
		// TODO: Generalize handling of log groups.
		for _, stream := range processLogStreams {
			logStreams = append(logStreams, fmt.Sprintf("%s:%s", group, stream))
		}
	}

	collector := log.CurrentLogCollector(ctx)
	collectorOutput, err := collector.GetEvents(ctx, &logd.GetEventsInput{
		Logs:   logStreams,
		Before: input.Before,
		After:  input.After,
	})
	if err != nil {
		return nil, err
	}
	output := api.GetEventsOutput{
		Events: make([]api.Event, len(collectorOutput.Events)),
		Cursor: collectorOutput.Cursor,
	}
	for i, collectorEvent := range collectorOutput.Events {
		output.Events[i] = api.Event{
			ID:        collectorEvent.ID,
			Log:       collectorEvent.Log,
			Timestamp: collectorEvent.Timestamp,
			Message:   collectorEvent.Message,
		}
	}
	return &output, nil
}

func (ws *Workspace) Start(ctx context.Context, input *api.StartInput) (*api.StartOutput, error) {
	id, err := ws.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving ref: %w", err)
	}

	components, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
		IDs:         []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("fetching component state: %w", err)
	}
	if len(components.Components) != 1 {
		return nil, fmt.Errorf("no state for component: %s", id)
	}
	component := components.Components[0]

	provider := ws.resolveProvider(ctx, component.Type)
	providerOutput, err := provider.Start(ctx, &core.StartInput{
		ID:    id,
		Spec:  component.Spec,
		State: component.State,
	})
	if err != nil {
		return nil, err
	}

	if _, err := ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:    id,
		State: providerOutput.State,
	}); err != nil {
		return nil, fmt.Errorf("updating component state: %w", err)
	}

	return &api.StartOutput{}, nil
}

func (ws *Workspace) Stop(ctx context.Context, input *api.StopInput) (*api.StopOutput, error) {
	id, err := ws.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving ref: %w", err)
	}

	components, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
		IDs:         []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("fetching component state: %w", err)
	}
	if len(components.Components) != 1 {
		return nil, fmt.Errorf("no state for component: %s", id)
	}
	component := components.Components[0]

	provider := ws.resolveProvider(ctx, component.Type)
	providerOutput, err := provider.Stop(ctx, &core.StopInput{
		ID:    id,
		Spec:  component.Spec,
		State: component.State,
	})
	if err != nil {
		return nil, err
	}

	if _, err := ws.Store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:    id,
		State: providerOutput.State,
	}); err != nil {
		return nil, fmt.Errorf("updating component state: %w", err)
	}

	return &api.StopOutput{}, nil
}

func (ws *Workspace) DescribeProcesses(ctx context.Context, input *api.DescribeProcessesInput) (*api.DescribeProcessesOutput, error) {
	components, err := ws.Store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		WorkspaceID: ws.ID,
		// TODO: Filter by type.
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	output := api.DescribeProcessesOutput{
		Processes: make([]api.ProcessDescription, 0, len(components.Components)),
	}
	for _, component := range components.Components {
		if component.Type == "process" {
			// XXX Do not utilize internal knowledge of process state.
			var state struct {
				Pid int `json:"pid"`
			}
			if err := jsonutil.UnmarshalString(component.State, &state); err != nil {
				// TODO: log error.
				fmt.Printf("unmarshalling process state: %v\n", err)
				continue
			}
			running := state.Pid != 0
			process := api.ProcessDescription{
				ID:      component.ID,
				Name:    component.Name,
				Running: running,
			}
			output.Processes = append(output.Processes, process)
		}
	}
	return &output, nil
}
