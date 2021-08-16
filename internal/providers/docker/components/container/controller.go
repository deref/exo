package container

import (
	"fmt"

	"github.com/deref/exo/internal/providers/docker"
	"github.com/deref/exo/internal/util/jsonutil"
)

func (c *Container) InitResource(componentID, spec, state string) error {
	c.ComponentID = componentID
	if err := docker.LoadSpec(spec, &c.Spec, c.ComponentBase.WorkspaceEnvironment); err != nil {
		return fmt.Errorf("loading spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(state, &c.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (c *Container) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(c.State)
}
