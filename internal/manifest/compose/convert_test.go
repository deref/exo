package compose_test

import (
	"strings"
	"testing"

	"github.com/deref/exo/internal/manifest"
	"github.com/deref/exo/internal/manifest/compose"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Skipf("Skipping convert test until we decide what strategy to use for converting compose files to component specs.")

	projectName := "testproj"
	defaultNetwork := manifest.Component{
		Name: "default",
		Type: "network",
		Spec: `driver: bridge
name: testproj_default
`,
	}

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
						defaultNetwork,
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
							DependsOn: []string{"default"},
						},
					},
				},
			},
		},

		{
			name: "named networks",
			in: `
services:
  proxy:
    image: nginx
    networks:
    - backend
    - frontend
  srv:
    image: myapp
    networks:
    - backend
networks:
  frontend:
  backend:
`,
			expected: manifest.LoadResult{
				Manifest: &manifest.Manifest{
					Components: []manifest.Component{
						defaultNetwork,
						{
							Name: "frontend",
							Type: "network",
							Spec: `driver: bridge
name: testproj_frontend
`,
						},
						{
							Name: "backend",
							Type: "network",
							Spec: `driver: bridge
name: testproj_backend
`,
						},
						{
							Name: "proxy",
							Type: "container",
							Spec: `container_name: testproj_proxy_1
image: nginx
networks:
- testproj_backend
- testproj_frontend
`,
							DependsOn: []string{"backend", "frontend"},
						},
						{
							Name: "srv",
							Type: "container",
							Spec: `container_name: testproj_srv_1
image: myapp
networks:
- testproj_backend
`,
							DependsOn: []string{"backend"},
						},
					},
				},
			},
		},
	}

	loader := compose.Loader{ProjectName: projectName}
	for _, testCase := range testCases {
		name := testCase.name
		in := strings.NewReader(testCase.in)
		expected := testCase.expected
		t.Run(name, func(t *testing.T) {
			out := loader.Load(in)

			assert.ElementsMatch(t, expected.Warnings, out.Warnings)
			if len(expected.Manifest.Components) > 0 {
				assert.ElementsMatch(t, expected.Manifest.Components, out.Manifest.Components)
			}
		})
	}
}
