package featureflag_test

import (
	"context"
	"testing"

	"github.com/deref/exo/internal/featureflag"
	"github.com/stretchr/testify/assert"
)

func TestIsEnabled(t *testing.T) {
	ctx := context.Background()
	fflags := featureflag.NewStaticFeatureFlags(map[string]bool{
		"a": true,
		"b": false,
	})

	assert.True(t, featureflag.IsEnabled(ctx, fflags, "a"))
	assert.False(t, featureflag.IsEnabled(ctx, fflags, "b"))
	assert.False(t, featureflag.IsEnabled(ctx, fflags, "c"))
}
