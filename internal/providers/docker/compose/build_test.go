package compose

import (
	"testing"
)

func TestBuildYAML(t *testing.T) {
	testYAML(t, "short", `./context/path`, Build{
		IsShortForm: true,
		BuildLongForm: BuildLongForm{
			Context: "./context/path",
		},
	})
	testYAML(t, "long", `
context: .
args:
  - x=y
shm_size: 1024
`, Build{
		BuildLongForm: BuildLongForm{
			Context: ".",
			Args: Dictionary{
				Style: MapStyle,
				Items: []DictionaryItem{
					DictionaryItem{
						Style: MapStyle,
						Key:   "x",
						Value: "y",
					},
				},
			},
			ShmSize: Bytes{
				Quantity: 1024,
			},
		},
	})
}
