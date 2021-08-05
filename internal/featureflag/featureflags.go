package featureflag

import "context"

type FeatureFlags interface {
	// CheckEnabled takes a list of flags and returns a parallel list
	// of booleans indicating whether the `i`th flag is enabled.
	CheckEnabled(ctx context.Context, flags []string) ([]bool, error)
}

func IsEnabled(ctx context.Context, ff FeatureFlags, flag string) bool {
	enabled, err := ff.CheckEnabled(ctx, []string{flag})
	if err != nil {
		// TODO: What is the right thing to do here?
		return false
	}

	return enabled[0]
}
