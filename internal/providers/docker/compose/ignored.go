package compose

import "gopkg.in/yaml.v3"

// This is a temporary placeholder for fields that we presently don't support,
// but are safe to ignore.
// TODO: Eliminate all usages of this with actual parsing logic.
type Ignored struct{}

func (_ *Ignored) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

func (_ *Ignored) Interpolate(env Environment) error {
	return nil
}
