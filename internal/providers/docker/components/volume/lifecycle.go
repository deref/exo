package volume

import (
	"context"
	"fmt"

	core "github.com/deref/exo/internal/core/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	dockerclient "github.com/docker/docker/client"
)

func (v *Volume) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	// See NOTE: [ADOPT COMPOSE RESOURCES].
	if existing, err := v.findExistingVolume(ctx); err != nil {
		return nil, fmt.Errorf("looking up existing volume: %w", err)
	} else if existing != nil {
		// TODO: Determine whether the existing volume is compatible with the spec.
		return &core.InitializeOutput{}, nil
	}

	labels := v.Spec.Labels.WithoutNils()
	for k, v := range v.GetExoLabels() {
		labels[k] = v
	}

	opts := volume.VolumeCreateBody{
		Driver:     v.Driver,
		DriverOpts: v.DriverOpts,
		Labels:     labels,
		Name:       v.Spec.Name,
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
		v.Logger.Infof("disposing volume not found: %q", v.VolumeName)
		err = nil
	}
	if err != nil {
		return nil, err
	}
	v.VolumeName = ""
	return &core.DisposeOutput{}, nil
}

func (v *Volume) findExistingVolume(ctx context.Context) (*types.Volume, error) {
	volume, err := v.Docker.VolumeInspect(ctx, v.Spec.Name)
	if err == nil {
		return &volume, nil
	}
	if dockerclient.IsErrNotFound(err) {
		return nil, nil
	}
	return nil, err
}
