package compose

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Dictionary struct {
	Style Style
	Items []DictionaryItem
}

type DictionaryItem struct {
	Style   Style
	Key     string
	Value   string
	NoValue bool
}

func (dict Dictionary) MarshalYAML() (interface{}, error) {
	if dict.Style == SeqStyle {
		return dict.Items, nil
	}
	node := &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, len(dict.Items)*2),
	}
	for i, item := range dict.Items {
		keyNode, valueNode := makeDictionaryItemNodes(item)
		node.Content[i*2+0] = &keyNode
		node.Content[i*2+1] = &valueNode
	}
	return node, nil
}

func (dict *Dictionary) UnmarshalYAML(node *yaml.Node) error {
	switch node.Tag {
	case "!!map":
		dict.Style = MapStyle
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
		dict.Style = SeqStyle
		return node.Decode(&dict.Items)
	default:
		return fmt.Errorf("expected !!seq or !!map, got %s", node.ShortTag())
	}
}

func (item DictionaryItem) MarshalYAML() (interface{}, error) {
	if item.Style == SeqStyle {
		if item.Value == "" {
			return item.Key, nil
		}
		return fmt.Sprintf("%s=%s", item.Key, item.Value), nil
	}
	keyNode, valueNode := makeDictionaryItemNodes(item)
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			&keyNode,
			&valueNode,
		},
	}, nil
}

func makeDictionaryItemNodes(item DictionaryItem) (keyNode, valueNode yaml.Node) {
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
		item.Style = SeqStyle
		parts := strings.SplitN(s, "=", 2)
		item.Key = parts[0]
		if len(parts) > 1 {
			item.Value = parts[1]
		} else {
			item.NoValue = true
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
	item.Style = MapStyle
	for k, v := range m {
		item.Key = k
		item.Value = v
	}
	return nil
}

func (dict Dictionary) Slice() []string {
	res := make([]string, len(dict.Items))
	for i, item := range dict.Items {
		if item.Value == "" {
			res[i] = item.Key
		} else {
			res[i] = fmt.Sprintf("%s=%s", item.Key, item.Value)
		}
	}
	return res
}

func (dict Dictionary) Map() map[string]string {
	m := make(map[string]string, len(dict.Items))
	for _, item := range dict.Items {
		m[item.Key] = item.Value
	}
	return m
}

func (dict Dictionary) MapOfPtr() map[string]*string {
	m := make(map[string]*string, len(dict.Items))
	for _, item := range dict.Items {
		if item.NoValue {
			m[item.Key] = nil
		} else {
			value := item.Value
			m[item.Key] = &value
		}
	}
	return m
}
