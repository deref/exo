package network

import (
	"fmt"

	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/util/jsonutil"
)

func (n *Network) InitResource(componentID, spec, state string) error {
	n.ComponentID = componentID
	if err := docker.LoadSpec(spec, &n.Spec, n.ComponentBase.WorkspaceEnvironment); err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(state, &n.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (n *Network) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(n.State)
}
