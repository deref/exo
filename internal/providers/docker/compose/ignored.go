package compose

// This is a temporary placeholder for fields that we presently don't support,
// but are safe to ignore.
// TODO: Eliminate all usages of this with actual parsing logic.
type Ignored struct{}

func (ignored *Ignored) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return nil
}
