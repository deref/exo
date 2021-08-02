package network

import (
	"context"

	core "github.com/deref/exo/core/api"
	"github.com/docker/docker/api/types"
)

func (n *Network) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	// Docker network names are non-unique aliases.
	// TODO: Use component name.
	name := n.ComponentID

	opts := types.NetworkCreate{
		// We don't care about duplicates, and it's best-effort checking only anyway.
		CheckDuplicate: false,
		Driver:         n.Driver,
		//Scope          string
		EnableIPv6: n.EnableIPv6,
		//IPAM           *network.IPAM
		Internal:   n.Internal,
		Attachable: n.Attachable,
		//Ingress        bool
		//ConfigOnly     bool
		//ConfigFrom     *network.ConfigReference
		//Options        map[string]string
		Labels: n.Labels.WithoutNils(),
	}
	createdBody, err := n.Docker.NetworkCreate(ctx, name, opts)
	if err != nil {
		return nil, err
	}

	n.NetworkID = createdBody.ID
	// TODO: Handle createdBody.Warnings.
	return &core.InitializeOutput{}, nil
}

func (n *Network) Update(context.Context, *core.UpdateInput) (*core.UpdateOutput, error) {
	panic("TODO: network update")
}

func (n *Network) Refresh(ctx context.Context, input *core.RefreshInput) (*core.RefreshOutput, error) {
	return &core.RefreshOutput{}, nil
}

func (n *Network) Dispose(ctx context.Context, input *core.DisposeInput) (*core.DisposeOutput, error) {
	if n.NetworkID == "" {
		return &core.DisposeOutput{}, nil
	}
	if err := n.Docker.NetworkRemove(ctx, n.NetworkID); err != nil {
		return nil, err
	}
	n.NetworkID = ""
	return &core.DisposeOutput{}, nil
}
