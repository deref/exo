package apigateway

import (
	"context"
	"fmt"

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
	webPortMapping, err := compose.ParsePortMappings(fmt.Sprintf("127.0.0.1:%d:8081", gatewaySpec.WebPort))
	if err != nil {
		ag.Logger.Infof("error: bad API gateway port mapping")
	}
	portMappings := append(apiPortMappings, webPortMapping...)

	token, err := ag.TokenClient.GetToken()
	if err != nil {
		return "", fmt.Errorf("getting token: %w", err)
	}

	return yamlutil.MustMarshalString(container.Spec{
		Image: compose.MakeString("f0ef8a799ce9"), // FIXME: host this docker image somewhere
		Environment: compose.Dictionary{Items: []compose.DictionaryItem{
			{Key: "EXO_URL", Value: "http://host.docker.internal:44643/_exo/"},
			{Key: "EXO_TOKEN", Value: token},
			{Key: "EXO_WORKSPACE_ID", Value: ag.WorkspaceID},
		}},
		ExtraHosts: compose.Strings{
			compose.MakeString("host.docker.internal:host-gateway"),
		},
		Ports: portMappings,
	}), nil
}

func (ag *APIGateway) Dependencies(ctx context.Context, input *core.DependenciesInput) (*core.DependenciesOutput, error) {
	return &core.DependenciesOutput{}, nil
}

func (ag *APIGateway) Initialize(ctx context.Context, input *core.InitializeInput) (output *core.InitializeOutput, err error) {
	var spec Spec
	if err := jsonutil.UnmarshalString(input.Spec, &spec); err != nil {
		return nil, fmt.Errorf("unmarshalling spec: %w", err)
	}

	containerSpec, err := ag.makeContainerSpec(spec)
	if err != nil {
		return nil, fmt.Errorf("making container spec: %w", err)
	}
	return ag.Container.Initialize(ctx, &core.InitializeInput{Spec: containerSpec})
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
	return ag.Container.Refresh(ctx, &core.RefreshInput{Spec: containerSpec})
}
