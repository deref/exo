package featureflag_test

import (
	"context"
	"testing"

	"github.com/deref/exo/internal/featureflag"
	"github.com/stretchr/testify/assert"
)

func TestStatic(t *testing.T) {
	ctx := context.Background()
	fflags := featureflag.NewStaticFeatureFlags(map[string]bool{
		"a": true,
		"b": false,
	})

	enabledFlags, err := fflags.CheckEnabled(ctx, []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Equal(t, []bool{true, false, false}, enabledFlags)
}
