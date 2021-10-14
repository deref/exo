package compose_test

import (
	"bytes"
	"testing"

	"github.com/deref/exo/internal/manifest/compose"
	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/deref/exo/internal/manifest/exohcl/testutil"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
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
				exo = "1.0"
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
				exo = "1.0"
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
			out, diags := converter.Convert([]byte(testCase.In))
			if len(diags) > 0 {
				assert.NoError(t, diags)
				return
			}
			var buf bytes.Buffer
			_, err := hclgen.WriteTo(&buf, &hcl.File{
				Body:  out.Body.(*hclsyntax.Body),
				Bytes: out.Bytes,
			})
			if assert.NoError(t, err) {
				// These tests are too brittle, as they are very sensitive to spacing and ordering.
				assert.Equal(t, testutil.CleanHCL([]byte(testCase.Expected)), testutil.CleanHCL(buf.Bytes()))
			}
		})
	}
}
