package network

import (
	"context"
	"errors"
	"fmt"

	core "github.com/deref/exo/internal/core/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (n *Network) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	if n.Name == "" {
		return nil, errors.New("Network must have a name")
	}

	if n.External {
		nets, err := n.Docker.NetworkList(ctx, types.NetworkListOptions{
			Filters: filters.NewArgs(filters.KeyValuePair{
				Key:   "name",
				Value: n.Name,
			}),
		})
		if err != nil {
			return nil, fmt.Errorf("listing networks: %w", err)
		}
		if len(nets) == 0 {
			return nil, fmt.Errorf("network %q not found", n.Name)
		}

		n.NetworkID = nets[0].ID

		return &core.InitializeOutput{}, nil
	}

	labels := n.Spec.Labels.WithoutNils()
	for k, v := range n.GetExoLabels() {
		labels[k] = v
	}

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
		Labels: labels,
	}
	createdBody, err := n.Docker.NetworkCreate(ctx, n.Name, opts)
	if err != nil {
		return nil, err
	}

	n.NetworkID = createdBody.ID
	// TODO: Handle createdBody.Warnings.
	return &core.InitializeOutput{}, nil
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
