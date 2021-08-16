package volume

import (
	"context"

	core "github.com/deref/exo/internal/core/api"
	"github.com/docker/docker/api/types/volume"
)

func (v *Volume) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	opts := volume.VolumeCreateBody{
		Driver: v.Driver.Value,
		Labels: v.Labels.WithoutNils(),
		Name:   v.Name.Value,
	}
	opts.DriverOpts = make(map[string]string)
	for k, v := range v.DriverOpts {
		opts.DriverOpts[k] = v.Value
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
	if err := v.Docker.VolumeRemove(ctx, v.VolumeName, force); err != nil {
		return nil, err
	}
	v.VolumeName = ""
	return &core.DisposeOutput{}, nil
}
