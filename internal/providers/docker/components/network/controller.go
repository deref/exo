package network

import (
	"fmt"

	"github.com/deref/util-go/jsonutil"
	"github.com/goccy/go-yaml"
)

func (n *Network) InitResource() error {
	if err := yaml.Unmarshal([]byte(n.ComponentSpec), &n.Spec); err != nil {
		return fmt.Errorf("unmarshalling spec: %w", err)
	}
	if err := jsonutil.UnmarshalString(n.ComponentState, &n.State); err != nil {
		return fmt.Errorf("unmarshalling state: %w", err)
	}
	return nil
}

func (n *Network) MarshalState() (state string, err error) {
	return jsonutil.MarshalString(n.State)
}
