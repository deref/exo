package exocue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractComponent(t *testing.T) {
	t.Skip("TODO")
}

func TestInjectComponent(t *testing.T) {
	t.Skip("TODO")
}

func TestComponentSpec(t *testing.T) {
	b := NewBuilder()
	b.AddManifest(`
		environment: {
			COMMON: "VAR"
		}
		components: {
			backend: $Daemon & {
				spec: {
					program: "./run-backend.sh"
					environment: {
						PORT: "1234"
					}
				}
			}
		}
	`)
	cfg := b.Build()
	spec := cfg.ComponentSpec("backend")
	assert.NoError(t, spec.Err())
	type Spec struct {
		Program     string
		Arguments   []string
		Environment map[string]interface{}
	}
	var actual Spec
	assert.NoError(t, spec.Decode(&actual))
	assert.Equal(t, Spec{
		Program:   "./run-backend.sh",
		Arguments: []string{},
		Environment: map[string]interface{}{
			"PORT":   "1234",
			"COMMON": "VAR",
		},
	}, actual)
}
