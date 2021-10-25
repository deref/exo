package network

import (
	"context"
	"errors"
	"fmt"

	core "github.com/deref/exo/internal/core/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

var _ core.Lifecycle = (*Network)(nil)

func (n *Network) Dependencies(ctx context.Context, input *core.DependenciesInput) (*core.DependenciesOutput, error) {
	return &core.DependenciesOutput{Components: []string{}}, nil
}

func (n *Network) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	var spec Spec
	if err := n.LoadSpec(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	if spec.Name.Value == "" {
		return nil, errors.New("Network must have a name")
	}

	existing, err := n.findExistingNetwork(ctx, spec.Name.Value)
	if err != nil {
		return nil, fmt.Errorf("looking up existing network: %w", err)
	}

	if spec.External.Value {
		if existing == nil {
			return nil, fmt.Errorf("network %q not found", spec.Name)
		}

		n.NetworkID = existing.ID

		return &core.InitializeOutput{}, nil
	}

	// NOTE [ADOPT COMPOSE RESOURCES]:
	// For networks and volumes, we look for existing resources that follow the docker-compose naming
	// convention and "adopt" them. That is, if they exist, we update the component ID to match the
	// ID of the existing resource, and skip initialization. This could be problematic if the user
	// happened to create a resource whose name conflicted but was not actually managed by compose.
	// Another potential difficulty is that we do not add our exo labels to an adopted resource because
	// that resouce would need to be recreated. This is not a problem for networks, but it could be
	// problematic with volumes that may contain data that the user expects to persist across an import.
	if existing != nil {
		// TODO: Determine whether the existing network is compatible with the spec.
		n.NetworkID = existing.ID
		return &core.InitializeOutput{}, nil
	}

	labels := spec.Labels.Map()
	for k, v := range n.GetExoLabels() {
		labels[k] = v
	}

	opts := types.NetworkCreate{
		// We don't care about duplicates, and it's best-effort checking only anyway.
		CheckDuplicate: false,
		Driver:         spec.Driver.Value,
		//Scope          string
		EnableIPv6: spec.EnableIPv6.Value,
		//IPAM           *network.IPAM
		Internal:   spec.Internal.Value,
		Attachable: spec.Attachable.Value,
		//Ingress        bool
		//ConfigOnly     bool
		//ConfigFrom     *network.ConfigReference
		//Options        map[string]string
		Labels: labels,
	}
	createdBody, err := n.Docker.NetworkCreate(ctx, spec.Name.Value, opts)
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
	err := n.Docker.NetworkRemove(ctx, n.NetworkID)
	if docker.IsErrNotFound(err) {
		n.Logger.Infof("network to be removed not found: %q", n.NetworkID)
		err = nil
	}
	if err != nil {
		return nil, err
	}
	n.NetworkID = ""
	return &core.DisposeOutput{}, nil
}

func (n *Network) findExistingNetwork(ctx context.Context, name string) (*types.NetworkResource, error) {
	nets, err := n.Docker.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: name,
		}),
	})
	if err != nil {
		return nil, err
	}

	switch len(nets) {
	case 0:
		return nil, nil
	case 1:
		return &nets[0], nil
	default:
		return nil, fmt.Errorf("expected 1 network but found %d", len(nets))
	}
}
