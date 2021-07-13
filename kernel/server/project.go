package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deref/exo/chrono"
	"github.com/deref/exo/components/invalid"
	"github.com/deref/exo/components/log"
	"github.com/deref/exo/components/process"
	"github.com/deref/exo/gensym"
	"github.com/deref/exo/kernel/api"
	"github.com/deref/exo/kernel/state"
	logcol "github.com/deref/exo/logcol/api"
)

type Project struct {
	ID string `json:"id"`
	// TODO: Path to root of directory.
}

func (proj *Project) Delete(ctx context.Context, input *api.DeleteInput) (*api.DeleteOutput, error) {
	store := state.CurrentStore(ctx)
	describeOutput, err := store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		ProjectID: proj.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	// TODO: Parallelism / bulk delete.
	for _, component := range describeOutput.Components {
		_, err := proj.DeleteComponent(ctx, &api.DeleteComponentInput{
			Ref: component.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("deleting %s: %w", component.Name, err)
		}
	}
	return &api.DeleteOutput{}, nil
}

func (proj *Project) Apply(ctx context.Context, input *api.ApplyInput) (*api.ApplyOutput, error) {
	panic("TODO: Apply")
}

func (proj *Project) Resolve(ctx context.Context, input *api.ResolveInput) (*api.ResolveOutput, error) {
	store := state.CurrentStore(ctx)
	storeOutput, err := store.Resolve(ctx, &state.ResolveInput{
		ProjectID: proj.ID,
		Refs:      input.Refs,
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

func (proj *Project) DescribeComponents(ctx context.Context, input *api.DescribeComponentsInput) (*api.DescribeComponentsOutput, error) {
	store := state.CurrentStore(ctx)
	stateOutput, err := store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		ProjectID: proj.ID,
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

func resolveLifecycle(typ string) api.Lifecycle {
	switch typ {
	case "process":
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		projectDir := wd                   // TODO: Get from component hierarchy.
		varDir := filepath.Join(wd, "var") // TODO: Get from exod config.
		return &process.Lifecycle{
			ProjectDir: projectDir,
			VarDir:     filepath.Join(varDir, "proc"),
		}
	default:
		return &invalid.Lifecycle{
			Err: fmt.Errorf("unsupported component type: %q", typ),
		}
	}
}

func (proj *Project) CreateComponent(ctx context.Context, input *api.CreateComponentInput) (*api.CreateComponentOutput, error) {
	if !IsValidName(input.Name) {
		return nil, fmt.Errorf("invalid name: %q", input.Name)
	}

	store := state.CurrentStore(ctx)

	id := gensym.Base32()

	if _, err := store.AddComponent(ctx, &state.AddComponentInput{
		ProjectID: "default",
		ID:        id,
		Name:      input.Name,
		Type:      input.Type,
		Spec:      input.Spec,
		Created:   chrono.NowString(ctx),
	}); err != nil {
		return nil, fmt.Errorf("adding component: %w", err)
	}

	lifecycle := resolveLifecycle(input.Type)
	output, err := lifecycle.Initialize(ctx, &api.InitializeInput{
		ID:   id,
		Spec: input.Spec,
	})
	if err != nil {
		return nil, err
	}

	if _, err := store.PatchComponent(ctx, &state.PatchComponentInput{
		ID:          id,
		State:       output.State,
		Initialized: chrono.NowString(ctx),
	}); err != nil {
		return nil, fmt.Errorf("modifying component after initialization: %w", err)
	}

	return &api.CreateComponentOutput{
		ID: id,
	}, nil
}

func IsValidName(name string) bool {
	return name != "" // XXX
}

func (proj *Project) UpdateComponent(ctx context.Context, input *api.UpdateComponentInput) (*api.UpdateComponentOutput, error) {
	panic("TODO: UpdateComponent")
}

func (proj *Project) RefreshComponent(ctx context.Context, input *api.RefreshComponentInput) (*api.RefreshComponentOutput, error) {
	panic("TODO: RefreshComponent")
}

func (proj *Project) DisposeComponent(ctx context.Context, input *api.DisposeComponentInput) (*api.DisposeComponentOutput, error) {
	id, err := proj.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving refs: %w", err)
	}
	err = proj.disposeComponent(ctx, id)
	return &api.DisposeComponentOutput{}, err
}

func (proj *Project) resolveRef(ctx context.Context, ref string) (string, error) {
	resolveOutput, err := proj.Resolve(ctx, &api.ResolveInput{Refs: []string{ref}})
	if err != nil {
		return "", err
	}
	id := resolveOutput.IDs[0]
	if id == nil {
		return "", nil
	}
	return *id, nil
}

func (proj *Project) disposeComponent(ctx context.Context, id string) error {
	store := state.CurrentStore(ctx)
	describeOutput, err := store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		ProjectID: proj.ID,
		IDs:       []string{id},
	})
	if err != nil {
		return fmt.Errorf("describing components: %w", err)
	}
	if len(describeOutput.Components) < 1 {
		return fmt.Errorf("no component %q", id)
	}
	component := describeOutput.Components[0]
	lifecycle := resolveLifecycle(component.Type)
	_, err = lifecycle.Dispose(ctx, &api.DisposeInput{
		ID:    id,
		State: component.State,
	})
	return err
}

func (proj *Project) DeleteComponent(ctx context.Context, input *api.DeleteComponentInput) (*api.DeleteComponentOutput, error) {
	id, err := proj.resolveRef(ctx, input.Ref)
	if err != nil {
		return nil, fmt.Errorf("resolving refs: %w", err)
	}
	if err := proj.disposeComponent(ctx, id); err != nil {
		return nil, fmt.Errorf("disposing: %w", err)
	}
	// TODO: Await disposal.
	store := state.CurrentStore(ctx)
	if _, err := store.RemoveComponent(ctx, &state.RemoveComponentInput{ID: id}); err != nil {
		return nil, fmt.Errorf("removing from state store: %w", err)
	}
	return &api.DeleteComponentOutput{}, nil
}

func (proj *Project) DescribeLogs(ctx context.Context, input *api.DescribeLogsInput) (*api.DescribeLogsOutput, error) {
	store := state.CurrentStore(ctx)
	components, err := store.DescribeComponents(ctx, &state.DescribeComponentsInput{
		ProjectID: proj.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("describing components: %w", err)
	}
	// TODO: When we have subcomponents, do a search for type=log.
	var logNames []string
	for _, component := range components.Components {
		if component.Type == "process" {
			for _, role := range []string{"out", "err"} {
				logNames = append(logNames, fmt.Sprintf("%s:%s", component.ID, role))
			}
		}
	}

	collector := log.CurrentLogCollector(ctx)
	collectorLogs, err := collector.DescribeLogs(ctx, &logcol.DescribeLogsInput{
		Names: logNames,
	})
	if err != nil {
		return nil, err
	}
	logs := make([]api.LogDescription, len(collectorLogs.Logs))
	for i, collectorLog := range collectorLogs.Logs {
		logs[i] = api.LogDescription{
			Name:        collectorLog.Name,
			LastEventAt: collectorLog.LastEventAt,
		}
	}
	return &api.DescribeLogsOutput{Logs: logs}, nil
}

func (proj *Project) GetEvents(ctx context.Context, input *api.GetEventsInput) (*api.GetEventsOutput, error) {
	collector := log.CurrentLogCollector(ctx)
	collectorEvents, err := collector.GetEvents(ctx, &logcol.GetEventsInput{
		Logs:   input.Logs,
		Before: input.Before,
		After:  input.After,
	})
	if err != nil {
		return nil, err
	}
	output := api.GetEventsOutput{
		Events: make([]api.Event, len(collectorEvents.Events)),
	}
	for i, collectorEvent := range collectorEvents.Events {
		output.Events[i] = api.Event{
			Log:       collectorEvent.Log,
			SID:       collectorEvent.SID,
			Timestamp: collectorEvent.Timestamp,
			Message:   collectorEvent.Message,
		}
	}
	return &output, nil
}
