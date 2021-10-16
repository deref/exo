package compose

import (
	"testing"
)

func TestBuildYAML(t *testing.T) {
	testYAML(t, "short", `./context/path`, Build{
		IsShortForm: true,
		BuildLongForm: BuildLongForm{
			Context: MakeString("./context/path"),
		},
	})
	testYAML(t, "long", `
context: .
args:
  - x=y
shm_size: 1024
`, Build{
		BuildLongForm: BuildLongForm{
			Context: MakeString("."),
			Args: Dictionary{
				Style: SeqStyle,
				Items: []DictionaryItem{
					{
						Style:  SeqStyle,
						String: MakeString("x=y"),
						Key:    "x",
						Value:  "y",
					},
				},
			},
			ShmSize: Bytes{
				String:   MakeInt(1024).String,
				Quantity: 1024,
			},
		},
	})
}
