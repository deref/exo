package volume

import (
	"context"
	"fmt"

	core "github.com/deref/exo/internal/core/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	dockerclient "github.com/docker/docker/client"
)

var _ core.Lifecycle = (*Volume)(nil)

func (v *Volume) Dependencies(ctx context.Context, input *core.DependenciesInput) (*core.DependenciesOutput, error) {
	return &core.DependenciesOutput{Components: []string{}}, nil
}

func (v *Volume) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	var spec Spec
	if err := v.LoadSpec(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	// See NOTE: [ADOPT COMPOSE RESOURCES].
	if existing, err := v.findExistingVolume(ctx, spec.Name.Value); err != nil {
		return nil, fmt.Errorf("looking up existing volume: %w", err)
	} else if existing != nil {
		// TODO: Determine whether the existing volume is compatible with the spec.
		return &core.InitializeOutput{}, nil
	}

	labels := spec.Labels.Map()
	for k, v := range v.GetExoLabels() {
		labels[k] = v
	}

	opts := volume.VolumeCreateBody{
		Driver:     spec.Driver.Value,
		DriverOpts: spec.DriverOpts.Map(),
		Labels:     labels,
		Name:       spec.Name.Value,
	}
	createdBody, err := v.Docker.VolumeCreate(ctx, opts)
	if err != nil {
		return nil, err
	}

	v.VolumeName = createdBody.Name
	// TODO: Capture more state from createdBody.
	return &core.InitializeOutput{}, nil
}

func (v *Volume) Refresh(ctx context.Context, input *core.RefreshInput) (*core.RefreshOutput, error) {
	return &core.RefreshOutput{}, nil
}

func (v *Volume) Dispose(ctx context.Context, input *core.DisposeInput) (*core.DisposeOutput, error) {
	if v.VolumeName == "" {
		return &core.DisposeOutput{}, nil
	}
	force := false
	err := v.Docker.VolumeRemove(ctx, v.VolumeName, force)
	if dockerclient.IsErrNotFound(err) {
		v.Logger.Infof("volume to be removed not found: %q", v.VolumeName)
		err = nil
	}
	if err != nil {
		return nil, err
	}
	v.VolumeName = ""
	return &core.DisposeOutput{}, nil
}

func (v *Volume) findExistingVolume(ctx context.Context, name string) (*types.Volume, error) {
	volume, err := v.Docker.VolumeInspect(ctx, name)
	if err == nil {
		return &volume, nil
	}
	if dockerclient.IsErrNotFound(err) {
		return nil, nil
	}
	return nil, err
}
