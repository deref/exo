package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/deref/exo/api"
	"github.com/deref/exo/chrono"
	"github.com/deref/exo/components/invalid"
	"github.com/deref/exo/components/process"
	"github.com/deref/exo/gensym"
	"github.com/deref/exo/state"
)

type Project struct {
	ID string `json:"id"`
	// TODO: Path to root of directory.
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

func (proj *Project) Apply(ctx context.Context, input *api.ApplyInput) (*api.ApplyOutput, error) {
	panic("TODO: Apply")
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
		ID:      id,
		Name:    input.Name,
		Type:    input.Type,
		Spec:    input.Spec,
		Created: chrono.NowString(ctx),
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
	panic("TODO: DisposeComponent")
}

func (proj *Project) DeleteComponent(ctx context.Context, input *api.DeleteComponentInput) (*api.DeleteComponentOutput, error) {
	panic("TODO: DeleteComponent")
}
