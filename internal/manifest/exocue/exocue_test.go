package exocue

import (
	"testing"

	"cuelang.org/go/cue"
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
	assert.NoError(t, cue.Value(spec).Err())
	type Spec struct {
		Program     string
		Arguments   []string
		Environment map[string]any
	}
	type Component struct {
		Name string
		Type string
		Spec Spec
	}
	var actual Component
	assert.NoError(t, cue.Value(spec).Decode(&actual))
	assert.Equal(t, Component{
		Name: "backend",
		Type: "daemon",
		Spec: Spec{
			Program:   "./run-backend.sh",
			Arguments: []string{},
			Environment: map[string]any{
				"PORT":          "1234",
				"COMMON":        "VAR",
				"EXO_COMPONENT": "backend",
			},
		},
	}, actual)
}
