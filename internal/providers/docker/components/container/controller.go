package container

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (c *Container) InitResource() error {
	if err := jsonutil.UnmarshalStringOrEmpty(c.ComponentState, &c.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (c *Container) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(c.State)
}
