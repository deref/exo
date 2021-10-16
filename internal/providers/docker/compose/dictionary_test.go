package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictionaryItemYAML(t *testing.T) {
	testYAML(t, "bare", `key`, DictionaryItem{
		Style:   SeqStyle,
		String:  MakeString("key"),
		Key:     "key",
		NoValue: true,
	})
	testYAML(t, "colon_empty", `key:`, DictionaryItem{
		Style: MapStyle,
		Key:   "key",
	})
	testYAML(t, "colon_value", `key: value`, DictionaryItem{
		Style:  MapStyle,
		String: MakeString("value"),
		Key:    "key",
		Value:  "value",
	})
	testYAML(t, "equal", `key=value`, DictionaryItem{
		Style:  SeqStyle,
		String: MakeString("key=value"),
		Key:    "key",
		Value:  "value",
	})
}

func TestDictionaryYAML(t *testing.T) {
	testYAML(t, "map", `
key: value
novalue:
`, Dictionary{
		Style: MapStyle,
		Items: []DictionaryItem{
			{
				Style:  MapStyle,
				String: MakeString("value"),
				Key:    "key",
				Value:  "value",
			},
			{
				Style: MapStyle,
				Key:   "novalue",
			},
		},
	})
	testYAML(t, "seq", `
- key=value
- novalue
`, Dictionary{
		Style: SeqStyle,
		Items: []DictionaryItem{
			{
				Style:  SeqStyle,
				String: MakeString("key=value"),
				Key:    "key",
				Value:  "value",
			},
			{
				Style:   SeqStyle,
				String:  MakeString("novalue"),
				Key:     "novalue",
				NoValue: true,
			},
		},
	})
}

func TestDictionarySlice(t *testing.T) {
	assert.Equal(t, []string{
		"novalue",
		"k=v",
	}, Dictionary{
		Items: []DictionaryItem{
			{
				Key: "novalue",
			},
			{
				Key:   "k",
				Value: "v",
			},
		},
	}.Slice())
}
