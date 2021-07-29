package container

import (
	"fmt"

	"github.com/deref/exo/util/jsonutil"
	"github.com/goccy/go-yaml"
)

func (c *Container) InitResource(componentID, spec, state string) error {
	c.ComponentID = componentID
	if err := yaml.Unmarshal([]byte(spec), &c.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(state, &c.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (c *Container) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(c.State)
}
