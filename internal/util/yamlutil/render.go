package yamlutil

import (
	yaml "gopkg.in/yaml.v3"
)

// MergeFormatting merges the formatting of the original yaml into the updated
// version. At the moment it's not tremendously sophisticated and just preserves
// comments.
func MergeFormatting(original *yaml.Node, updated *yaml.Node) error {
	if original == nil {
		return nil
	}

	// TODO: preserve order of maps
	if updated.FootComment == "" {
		updated.FootComment = original.FootComment
	}
	if updated.HeadComment == "" {
		updated.HeadComment = original.HeadComment
	}
	if updated.LineComment == "" {
		updated.LineComment = original.LineComment
	}

	if original.Kind != updated.Kind {
		return nil
	}

	for i := range updated.Content {
		if len(original.Content) > i {
			if err := MergeFormatting(original.Content[i], updated.Content[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
