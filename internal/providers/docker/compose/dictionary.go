package compose

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type DictionarySyntax rune

const (
	DictionarySyntaxUnknown DictionarySyntax = 0
	DictionarySyntaxMap                      = 'M'
	DictionarySyntaxArray                    = 'A'
)

type Dictionary struct {
	Syntax DictionarySyntax
	Items  []DictionaryItem
}

type DictionaryItem struct {
	Syntax DictionarySyntax
	Key    string
	Value  string
}

func (dict Dictionary) MarshalYAML() (interface{}, error) {
	if dict.Syntax == DictionarySyntaxArray {
		return dict.Items, nil
	}
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, len(dict.Items)*2),
	}
	for i, item := range dict.Items {
		keyNode, valueNode := mapItemNodes(item)
		node.Content[i*2+0] = &keyNode
		node.Content[i*2+1] = &valueNode
	}
	return node, nil
}

func (dict *Dictionary) UnmarshalYAML(node *yaml.Node) error {
	switch node.Tag {
	case "!!map":
		dict.Syntax = DictionarySyntaxMap
		n := len(node.Content) / 2
		dict.Items = make([]DictionaryItem, n)
		for i := 0; i < n; i++ {
			if err := dict.Items[i].UnmarshalYAML(&yaml.Node{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					node.Content[i*2+0],
					node.Content[i*2+1],
				},
			}); err != nil {
				return err
			}
		}
		return nil
	case "!!seq":
		dict.Syntax = DictionarySyntaxArray
		return node.Decode(&dict.Items)
	default:
		return fmt.Errorf("expected !!seq or !!map, got %s", node.ShortTag())
	}
}

func (item DictionaryItem) MarshalYAML() (interface{}, error) {
	if item.Syntax == DictionarySyntaxArray {
		if item.Value == "" {
			return item.Key, nil
		}
		return fmt.Sprintf("%s=%s", item.Key, item.Value), nil
	}
	keyNode, valueNode := mapItemNodes(item)
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			&keyNode,
			&valueNode,
		},
	}, nil
}

func mapItemNodes(item DictionaryItem) (keyNode, valueNode yaml.Node) {
	if err := keyNode.Encode(item.Key); err != nil {
		panic(err)
	}
	if item.Value == "" {
		valueNode.Kind = yaml.ScalarNode
	} else {
		if err := valueNode.Encode(item.Value); err != nil {
			panic(err)
		}
	}
	return
}

func (item *DictionaryItem) UnmarshalYAML(node *yaml.Node) error {
	var s string
	err := node.Decode(&s)
	if err == nil {
		item.Syntax = DictionarySyntaxArray
		parts := strings.SplitN(s, "=", 2)
		item.Key = parts[0]
		if len(parts) > 1 {
			item.Value = parts[1]
		}
		return nil
	}

	var m map[string]string
	if err := node.Decode(&m); err != nil {
		return err
	}
	if len(m) != 1 {
		return errors.New("expected single mapping")
	}
	item.Syntax = DictionarySyntaxMap
	for k, v := range m {
		item.Key = k
		item.Value = v
	}
	return nil
}

/*

func (dict Dictionary) Slice() []string {
	m := map[string]*string(dict)
	res := make([]string, len(m))
	i := 0
	for k, v := range m {
		if v == nil {
			res[i] = k
		} else {
			res[i] = fmt.Sprintf("%s=%s", k, *v)
		}
		i++
	}
	sort.Strings(res)
	return res
}

func (dict Dictionary) WithoutNils() map[string]string {
	m := make(map[string]string, len(dict))
	for k, v := range dict {
		if v == nil {
			continue
		}
		m[k] = *v
	}
	return m
}

*/
