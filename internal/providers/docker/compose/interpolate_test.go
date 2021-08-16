package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpolate(t *testing.T) {
	tests := []struct {
		Expected    Compose
		Subject     Compose
		Environment map[string]string
	}{
		{
			Compose{
				Version: String{
					Expression: "${VERSION}",
					Value:      "123",
				},
			},
			Compose{
				Version: String{
					Expression: "${VERSION}",
				},
			},
			map[string]string{
				"VERSION": "123",
			},
		},
	}
	for _, test := range tests {
		err := Interpolate(&test.Subject, MapEnvironment(test.Environment))
		assert.NoError(t, err)
		assert.Equal(t, test.Expected, test.Subject)
	}
}
