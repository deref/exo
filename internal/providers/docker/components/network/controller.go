package network

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
)

func (n *Network) InitResource() error {
	if err := jsonutil.UnmarshalStringOrEmpty(n.ComponentState, &n.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (n *Network) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(n.State)
}
