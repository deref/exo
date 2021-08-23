package compose_test

import (
	"strings"
	"testing"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/compose"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected manifest.LoadResult
	}{
		{
			name: "basic service",
			in: `
services:
    web:
        image: "nodejs:14"
        volumes: ['./src:/srv']
        command: node /srv/index.js
`,
			expected: manifest.LoadResult{
				Manifest: &manifest.Manifest{
					Components: []manifest.Component{
						{
							Name: "default",
							Type: "network",
							Spec: `driver: bridge
name: testproj_default
`,
						},
						{
							Name: "web",
							Type: "container",
							Spec: `command: node /srv/index.js
container_name: testproj_web_1
image: nodejs:14
networks:
- testproj_default
volumes:
- ./src:/srv
`,
						},
					},
				},
			},
		},
	}

	loader := compose.Loader{ProjectName: "testproj"}
	for _, testCase := range testCases {
		name := testCase.name
		in := strings.NewReader(testCase.in)
		expected := testCase.expected
		t.Run(name, func(t *testing.T) {
			out := loader.Load(in)
			assert.Equal(t, expected, out)
		})
	}
}
