package procfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrganize(t *testing.T) {
	tests := []struct {
		Processes []Process
		Expected  []Process
	}{
		{
			[]Process{},
			[]Process{},
		},

		// No assigned ports; alphabetize.
		{
			[]Process{
				{Name: "z"},
				{Name: "a"},
			},
			[]Process{
				{Name: "a"},
				{Name: "z"},
			},
		},

		// Correctly assigned port takes precedence over alphabetization.
		{
			[]Process{
				{Name: "z", Environment: map[string]string{"PORT": "5000"}},
				{Name: "a"},
			},
			[]Process{
				{Name: "z", Environment: map[string]string{}},
				{Name: "a"},
			},
		},

		// Unaligned port in range.
		{
			[]Process{
				{Name: "z", Environment: map[string]string{"PORT": "5001"}},
				{Name: "a"},
			},
			[]Process{
				{Name: "a"},
				{Name: "z", Environment: map[string]string{"PORT": "5001"}},
			},
		},

		// Range, gapless.
		{
			[]Process{
				{Name: "c", Environment: map[string]string{"PORT": "5200"}},
				{Name: "b", Environment: map[string]string{"PORT": "5100"}},
				{Name: "a", Environment: map[string]string{"PORT": "5000"}},
			},
			[]Process{
				{Name: "a", Environment: map[string]string{}},
				{Name: "b", Environment: map[string]string{}},
				{Name: "c", Environment: map[string]string{}},
			},
		},

		// Gap.
		{
			[]Process{
				{Name: "a", Environment: map[string]string{"PORT": "5000"}},
				{Name: "c", Environment: map[string]string{"PORT": "5200"}},
				{Name: "z"},
				{Name: "b", Environment: map[string]string{"PORT": "5100"}},
			},
			[]Process{
				{Name: "a", Environment: map[string]string{}},
				{Name: "b", Environment: map[string]string{}},
				{Name: "c", Environment: map[string]string{}},
				{Name: "z"},
			},
		},
	}
	for _, test := range tests {
		Organize(&test.Processes)
		assert.Equal(t, test.Expected, test.Processes)
	}
}
