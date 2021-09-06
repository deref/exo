package procfile_test

import (
	"bytes"
	"testing"

	"github.com/deref/exo/internal/manifest/procfile"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		name     string
		in       []procfile.Process
		expected string
	}{
		{
			name: "single process",
			in: []procfile.Process{
				{
					Name:      "server",
					Program:   "npm",
					Arguments: []string{"run"},
					Environment: map[string]string{
						"NODE_ENV": "development",
					},
				},
			},
			expected: "server: NODE_ENV=development npm run\n",
		},

		{
			name: "multiple processes",
			in: []procfile.Process{
				{
					Name:      "server",
					Program:   "npm",
					Arguments: []string{"run"},
					Environment: map[string]string{
						"NODE_ENV": "development",
					},
				},
				{
					Name:        "ui",
					Program:     "docker",
					Arguments:   []string{"run", "--rm", "--name=nginx", "-v", "config/nginx.conf:/etc/nginx/nginx.conf:ro", "-d", "nginx"},
					Environment: map[string]string{},
				},
			},
			expected: `server: NODE_ENV=development npm run
ui: docker run --rm --name=nginx -v config/nginx.conf:/etc/nginx/nginx.conf:ro -d nginx
`,
		},

		{
			name: "escaped content",
			in: []procfile.Process{
				{
					Name:      "weird",
					Program:   "foo'",
					Arguments: []string{"multiple words", "Hello\nWorld", `"quotes"`, `'other quotes'`},
					Environment: map[string]string{
						"SOME_ARG":  "ignore=other=equals",
						"OTHER_ARG": "'\"",
					},
				},
			},
			expected: `weird: OTHER_ARG=''"'"'"' SOME_ARG=ignore=other=equals 'foo'"'"'' 'multiple words' 'Hello
World' '"quotes"' ''"'"'other quotes'"'"''
`,
		},
	}

	for _, testCase := range testCases {
		in := testCase.in
		expected := testCase.expected
		var buf bytes.Buffer
		err := procfile.Generate(&buf, in)
		if !assert.NoError(t, err) {
			break
		}

		assert.Equal(t, expected, buf.String())
	}
}
