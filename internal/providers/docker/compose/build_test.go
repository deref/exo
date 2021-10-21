package compose

import (
	"testing"
)

func TestBuildYAML(t *testing.T) {
	testYAML(t, "short", `./context/path`, Build{
		ShortForm: MakeString("./context/path"),
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
	assertInterpolated(t, map[string]string{"context": "./path"}, `${context}`, Build{
		ShortForm: MakeString("${context}").WithValue("./path"),
		BuildLongForm: BuildLongForm{
			Context: MakeString("${context}").WithValue("./path"),
		},
	})
	assertInterpolated(t, map[string]string{"x": "1"}, `
dockerfile: foo${x}bar
`, Build{
		BuildLongForm: BuildLongForm{
			Dockerfile: MakeString("foo${x}bar").WithValue("foo1bar"),
		},
	})
}
