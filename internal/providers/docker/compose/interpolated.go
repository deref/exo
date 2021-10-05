package compose

import (
	"fmt"

	"github.com/deref/exo/internal/providers/docker/compose/interpolate"
	"github.com/goccy/go-yaml"
)

type Environment = interpolate.Environment
type MapEnvironment = interpolate.MapEnvironment

// Interpolated unmarshalls into a generic structure, performs interpolation
// using the Environment set on the Interpolated instance, and then unmarshalls
// into Value.
type Interpolated struct {
	Environment Environment
	Value       interface{}
}

func (i *Interpolated) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Decode into generic data types.
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	// Interpolate and then convert back to yaml.
	// Ideally, we could bypass the extra marshal/unmarshal pair, but the
	// interface unmarshalling internals are not exposed from the yaml library.
	if err := interpolate.Interpolate(raw, i.Environment); err != nil {
		return fmt.Errorf("interpolating: %w", err)
	}
	interpolatedBytes, err := yaml.Marshal(raw)
	if err != nil {
		// Should be unreachable, but potentially possible if there is some
		// weird structure that can be unmarshalled but not marshalled.
		return fmt.Errorf("intermediate remarshaling: %w", err)
	}

	// Decode into stronger data types.
	if err := yaml.Unmarshal(interpolatedBytes, i.Value); err != nil {
		return err
	}
	return nil
}
