package apigateway

import (
	"context"
	"fmt"
	"strconv"

	core "github.com/deref/exo/internal/core/api"
	"github.com/deref/exo/internal/providers/docker/components/container"
	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/yamlutil"
)

func (ag APIGateway) makeContainerSpec(gatewaySpec Spec) (string, error) {
	apiPortMappings, err := compose.ParsePortMappings(fmt.Sprintf("127.0.0.1:%d:8080", gatewaySpec.APIPort))
	if err != nil {
		ag.Logger.Infof("error: bad API gateway port mapping")
	}
	webPortMapping, err := compose.ParsePortMappings(fmt.Sprintf("127.0.0.1:%d:8082", gatewaySpec.WebPort))
	if err != nil {
		ag.Logger.Infof("error: bad web port mapping")
	}
	portMappings := append(apiPortMappings, webPortMapping...)

	token, err := ag.TokenClient.GetToken()
	if err != nil {
		return "", fmt.Errorf("getting token: %w", err)
	}

	return yamlutil.MustMarshalString(container.Spec{
		Image: compose.MakeString("ghcr.io/deref/exo-mitm:latest"),
		Environment: compose.Dictionary{Items: []compose.DictionaryItem{
			{Key: "EXO_URL", Value: fmt.Sprintf("http://host.docker.internal:%d/_exo/", ag.HTTPPort)},
			{Key: "EXO_TOKEN", Value: token},
			{Key: "EXO_WORKSPACE_ID", Value: ag.WorkspaceID},
		}},
		ExtraHosts: compose.Strings{
			compose.MakeString("host.docker.internal:host-gateway"),
		},
		Ports: portMappings,
	}), nil
}

func (ag *APIGateway) UpdatePorts(ctx context.Context) error {
	container, err := ag.Docker.ContainerInspect(ctx, ag.State.ContainerID)
	if err != nil {
		return fmt.Errorf("inspecting container: %w", err)
	}

	apiPortMappings, okay := container.NetworkSettings.Ports["8080/tcp"]
	if !okay || len(apiPortMappings) == 0 {
		return fmt.Errorf("API port mapping not found")
	}
	apiPort, err := strconv.Atoi(apiPortMappings[0].HostPort)
	if err != nil {
		return fmt.Errorf("converting API port to int: %w", err)
	}

	webPortMappings, okay := container.NetworkSettings.Ports["8082/tcp"]
	if !okay || len(webPortMappings) == 0 {
		return fmt.Errorf("host port mapping not found")
	}
	webPort, err := strconv.Atoi(webPortMappings[0].HostPort)
	if err != nil {
		return fmt.Errorf("converting host port to int: %w", err)
	}

	ag.State.APIPort = apiPort
	ag.State.WebPort = webPort
	return nil
}

func (ag *APIGateway) Dependencies(ctx context.Context, input *core.DependenciesInput) (*core.DependenciesOutput, error) {
	return &core.DependenciesOutput{}, nil
}

func (ag *APIGateway) Initialize(ctx context.Context, input *core.InitializeInput) (*core.InitializeOutput, error) {
	var spec Spec
	if err := jsonutil.UnmarshalString(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
	}

	containerSpec, err := ag.makeContainerSpec(spec)
	if err != nil {
		return nil, fmt.Errorf("making container spec: %w", err)
	}

	output, err := ag.Container.Initialize(ctx, &core.InitializeInput{Spec: containerSpec})
	if err != nil {
		return nil, err
	}
	ag.State.State = ag.Container.State

	if err := ag.UpdatePorts(ctx); err != nil {
		return nil, fmt.Errorf("updating ports: %w", err)
	}

	return output, nil
}

func (ag *APIGateway) Refresh(ctx context.Context, input *core.RefreshInput) (*core.RefreshOutput, error) {
	var spec Spec
	if err := jsonutil.UnmarshalString(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
	}

	containerSpec, err := ag.makeContainerSpec(spec)
	if err != nil {
		return nil, fmt.Errorf("making container spec: %w", err)
	}
	output, err := ag.Container.Refresh(ctx, &core.RefreshInput{Spec: containerSpec})
	if err != nil {
		return nil, err
	}
	ag.State.State = ag.Container.State

	if err := ag.UpdatePorts(ctx); err != nil {
		return nil, fmt.Errorf("updating ports: %w", err)
	}

	return output, nil
}
