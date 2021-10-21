package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/deref/exo/internal/core/api"
)

func (ws *Workspace) RenderDependencies(ctx context.Context, input *api.RenderDependenciesInput) (*api.RenderDependenciesOutput, error) {
	sb := &strings.Builder{}
	sb.WriteString("digraph {\n")

	describe, err := ws.DescribeComponents(ctx, &api.DescribeComponentsInput{})
	if err != nil {
		return nil, fmt.Errorf("describing components")
	}
	for _, component := range describe.Components {
		var deps *api.DependenciesOutput
		if err := ws.query(ctx, component, &deps, &api.DependenciesInput{
			Spec: component.Spec,
		}); err != nil {
			return nil, fmt.Errorf("getting %q dependencies: %w", component.Name, err)
		}
		fmt.Fprintf(sb, "%q;\n", component.Name)
		for _, dep := range deps.Components {
			fmt.Fprintf(sb, "%q -> %q;\n", component.Name, dep)
		}
	}

	sb.WriteString("}\n")
	return &api.RenderDependenciesOutput{
		Dot: sb.String(),
	}, nil
}
