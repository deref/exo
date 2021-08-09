package network

import (
	"fmt"

	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/goccy/go-yaml"
)

func (n *Network) InitResource(spec, state string) error {
	if err := yaml.Unmarshal([]byte(spec), &n.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(state, &n.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (n *Network) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(n.State)
}
