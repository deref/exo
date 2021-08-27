package compose_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/stretchr/testify/assert"
)

func TestParseService(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected compose.Service
	}{
		{
			name: "volumes - full syntax",
			in: `volumes:
- type: volume
  source: mydata
  target: /data
  read_only: true
  volume:
    nocopy: true
- type: bind
  source: /path/a
  target: /path/b
  bind:
    propagation: rshared
    create_host_path: true
- type: tmpfs
  target: /data/buffer
  tmpfs:
    size: 208666624`,
			expected: compose.Service{
				Volumes: []compose.VolumeMount{
					{
						Type:     "volume",
						Source:   "mydata",
						Target:   "/data",
						ReadOnly: true,
						Volume: &compose.VolumeOptions{
							Nocopy: true,
						},
					},
					{
						Type:   "bind",
						Source: "/path/a",
						Target: "/path/b",
						Bind: &compose.BindOptions{
							Propagation:    "rshared",
							CreateHostPath: true,
						},
					},
					{
						Type:   "tmpfs",
						Target: "/data/buffer",
						Tmpfs: &compose.TmpfsOptions{
							Size: 208666624,
						},
					},
				},
			},
		},

		{
			name: "volumes: short syntax",
			in: `volumes:
- /var/myapp
- './data:/data'
- "/home/fred/.ssh:/root/.ssh:ro"
- '~/util:/usr/bin/util:rw'
- my-log-volume:/var/log/xyzzy`,
			expected: compose.Service{
				Volumes: []compose.VolumeMount{
					{
						Type:   "volume",
						Target: "/var/myapp",
					},
					{
						Type:   "bind",
						Source: "./data",
						Target: "/data",
						Bind: &compose.BindOptions{
							CreateHostPath: true,
						},
					},
					{
						Type:     "bind",
						Source:   "/home/fred/.ssh",
						Target:   "/root/.ssh",
						ReadOnly: true,
						Bind: &compose.BindOptions{
							CreateHostPath: true,
						},
					},
					{
						Type:   "bind",
						Source: "~/util",
						Target: "/usr/bin/util",
						Bind: &compose.BindOptions{
							CreateHostPath: true,
						},
					},
					{
						Type:   "volume",
						Source: "my-log-volume",
						Target: "/var/log/xyzzy",
					},
				},
			},
		},

		{
			name: "service dependencies - short syntax",
			in: `depends_on:
- db
- messages`,
			expected: compose.Service{
				DependsOn: compose.ServiceDependencies{
					IsShortSyntax: true,
					Services: []compose.ServiceDependency{
						{
							Service:   "db",
							Condition: "service_started",
						},
						{
							Service:   "messages",
							Condition: "service_started",
						},
					},
				},
			},
		},

		{
			name: "service dependencies - extended syntax",
			in: `depends_on:
  db:
  messages:
    condition: service_healthy`,
			expected: compose.Service{
				DependsOn: compose.ServiceDependencies{
					Services: []compose.ServiceDependency{
						{
							Service:   "db",
							Condition: "service_started",
						},
						{
							Service:   "messages",
							Condition: "service_healthy",
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		name := testCase.name
		in := testCase.in
		expected := testCase.expected
		t.Run(name, func(t *testing.T) {
			var content bytes.Buffer
			content.WriteString("services:\n  test-svc:\n")
			lines := strings.Split(in, "\n")
			for _, line := range lines {
				content.WriteString("    ")
				content.WriteString(line)
				content.WriteByte('\n')
			}

			comp, err := compose.Parse(&content)
			assert.NoError(t, err)

			svc := comp.Services["test-svc"]
			assert.Equal(t, expected, svc)
		})
	}
}
