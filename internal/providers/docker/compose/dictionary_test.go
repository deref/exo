package compose

import (
	"testing"
)

func TestDictionaryItemYAML(t *testing.T) {
	testYAML(t, "bare", `key`, DictionaryItem{
		Syntax: DictionarySyntaxArray,
		Key:    "key",
	})
	testYAML(t, "colon_empty", `key:`, DictionaryItem{
		Syntax: DictionarySyntaxMap,
		Key:    "key",
	})
	testYAML(t, "colon_value", `key: value`, DictionaryItem{
		Syntax: DictionarySyntaxMap,
		Key:    "key",
		Value:  "value",
	})
	testYAML(t, "equal", `key=value`, DictionaryItem{
		Syntax: DictionarySyntaxArray,
		Key:    "key",
		Value:  "value",
	})
}

func TestDictionaryYAML(t *testing.T) {
	testYAML(t, "map", `
key: value
novalue:
`, Dictionary{
		Syntax: DictionarySyntaxMap,
		Items: []DictionaryItem{
			{
				Syntax: DictionarySyntaxMap,
				Key:    "key",
				Value:  "value",
			},
			{
				Syntax: DictionarySyntaxMap,
				Key:    "novalue",
			},
		},
	})
	testYAML(t, "slice", `
- key=value
- novalue
`, Dictionary{
		Syntax: DictionarySyntaxArray,
		Items: []DictionaryItem{
			{
				Syntax: DictionarySyntaxArray,
				Key:    "key",
				Value:  "value",
			},
			{
				Syntax: DictionarySyntaxArray,
				Key:    "novalue",
			},
		},
	})
}
