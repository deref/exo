package featureflag

import "context"

func NewStaticFeatureFlags(flags map[string]bool) *StaticFeatureFlags {
	return &StaticFeatureFlags{
		flags: flags,
	}
}

type StaticFeatureFlags struct {
	flags map[string]bool
}

func (f *StaticFeatureFlags) CheckEnabled(ctx context.Context, flags []string) ([]bool, error) {
	out := make([]bool, len(flags))
	for idx, flag := range flags {
		if isSet, ok := f.flags[flag]; ok && isSet {
			out[idx] = true
		}
	}

	return out, nil
}
