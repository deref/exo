package compose

import (
	"testing"
)

func TestDictionaryItemYAML(t *testing.T) {
	testYAML(t, "bare", `key`, DictionaryItem{
		Style: SeqStyle,
		Key:   "key",
	})
	testYAML(t, "colon_empty", `key:`, DictionaryItem{
		Style: MapStyle,
		Key:   "key",
	})
	testYAML(t, "colon_value", `key: value`, DictionaryItem{
		Style: MapStyle,
		Key:   "key",
		Value: "value",
	})
	testYAML(t, "equal", `key=value`, DictionaryItem{
		Style: SeqStyle,
		Key:   "key",
		Value: "value",
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
				Style: MapStyle,
				Key:   "key",
				Value: "value",
			},
			{
				Style: MapStyle,
				Key:   "novalue",
			},
		},
	})
	testYAML(t, "slice", `
- key=value
- novalue
`, Dictionary{
		Style: SeqStyle,
		Items: []DictionaryItem{
			{
				Style: SeqStyle,
				Key:   "key",
				Value: "value",
			},
			{
				Style: SeqStyle,
				Key:   "novalue",
			},
		},
	})
}
