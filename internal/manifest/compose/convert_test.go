package compose_test

import (
	"testing"

	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func TestConvert(t *testing.T) {
	// TODO: These tests are badly broken, since this entire approach is badly broken.
	// The following things that are being tested for here should not appear in converted manifests:
	// - prefixed container/volume/network names.
	// - docker compose labels
	// - dependencies that can be inferred from the spec body
	// XXX resume work here.

	projectName := "testproj"

	testCases := []struct {
		Name     string
		In       string
		Expected string
	}{
		{
			Name: "basic service",
			In: `
services:
    web:
        image: "nodejs:14"
        volumes: ['./src:/srv']
        command: node /srv/index.js
`,
			Expected: `
exo = "0.1"
components {
	network "default" {
		driver = "bridge"
		name = "testproj_default"
	}
	container "web" {
		command = "node /srv/index.js"
		container_name = "testproj_web_1"
		image = "nodejs:14"
		labels = { "com.docker.compose.project" = "testproj", "com.docker.compose.service" = "web" }
		networks = ["testproj_default"]
		volumes = ["./src:srv"]
		_ {
			depends_on = ["default"]
		}
	}
}`,
		},

		{
			Name: "named networks",
			In: `
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
			Expected: `
exo = "0.1"
components {
	network "frontend" {
		driver = "bridge"
		name = "testproj_frontend"
	}
	network "backend" {
		driver = "bridge"
		name = "testproj_backend"
	}
	network "default" {
		driver = "bridge"
		name = "testproj_default"
	}
	container "proxy" {
		container_name = "testproj_proxy_1"
		image = "nginx"
		labels = { "com.docker.compose.project" = "testproj", "com.docker.compose.service" = "proxy" }
		networks = [ "testproj_backend", "testproj_frontend" ]
		_ {
			depends_on = ["backend", "frontend"]
		}
	}
	container "srv" {
		container_name = "testproj_srv_1"
		image = "myapp"
		labels = { "com.docker.compose.project" = "testproj", "com.docker.compose.service" = "srv" }
		networks = ["testproj_backend"]
		_ {
			depends_on = ["backend"]
		}
	}
}`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			converter := &compose.Converter{ProjectName: projectName}
			actual, diags := converter.Convert([]byte(testCase.In))
			if len(diags) > 0 {
				t.Fatalf("error converting: %v", diags)
			}

			expected, diags := hclsyntax.ParseConfig([]byte(testCase.Expected), testCase.Name, hcl.InitialPos)
			if len(diags) > 0 {
				t.Fatalf("malformed test case: %v", diags)
			}

			if !hclgen.FileEquiv(expected, &hcl.File{Body: actual.Body}) {
				t.Errorf("hcl files inequivalent. expected:\n%s\nactual:\n%s",
					testCase.Expected,
					string(hclgen.FormatFile(actual)))
			}
		})
	}
}
