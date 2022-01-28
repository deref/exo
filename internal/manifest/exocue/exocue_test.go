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

func TestManifestToComponent(t *testing.T) {
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
	spec := cfg.Component("backend")
	assert.NoError(t, spec.Err())
	type Spec struct {
		Program     string
		Arguments   []string
		Environment map[string]interface{}
	}
	type Component struct {
		Name string
		Type string
		Spec Spec
	}
	var actual Component
	assert.NoError(t, spec.Decode(&actual))
	assert.Equal(t, Component{
		Name: "backend",
		Type: "daemon",
		Spec: Spec{
			Program:   "./run-backend.sh",
			Arguments: []string{},
			Environment: map[string]interface{}{
				"PORT":          "1234",
				"COMMON":        "VAR",
				"EXO_COMPONENT": "backend",
			},
		},
	}, actual)
}
