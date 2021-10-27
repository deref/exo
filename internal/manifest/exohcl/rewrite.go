package exohcl

import (
	"context"

	"github.com/deref/exo/internal/manifest/exohcl/hclgen"
	"github.com/hashicorp/hcl/v2"
)

type Rewrite interface {
	context.Context
	RewriteManifest(Rewrite, *Manifest) *hclgen.File
	RewriteEnvironment(Rewrite, *Environment) *hclgen.Block
	RewriteComponents(Rewrite, *ComponentSet) *hclgen.Block
	RewriteComponent(Rewrite, *Component) *hclgen.Block
}

func RewriteManifest(re Rewrite, m *Manifest) *hclgen.File {
	return re.RewriteManifest(re, m)
}

func RewriteEnvironment(re Rewrite, env *Environment) *hclgen.Block {
	return re.RewriteEnvironment(re, env)
}

func RewriteComponents(re Rewrite, cs *ComponentSet) *hclgen.Block {
	return re.RewriteComponents(re, cs)
}

func RewriteComponent(re Rewrite, c *Component) *hclgen.Block {
	return re.RewriteComponent(re, c)
}

type RewriteBase struct{}

func (_ RewriteBase) RewriteManifest(re Rewrite, m *Manifest) *hclgen.File {
	skip := make(map[*hcl.Block]bool)
	blocks := make([]*hclgen.Block, 0, len(m.Content.Blocks))

	env := NewEnvironment(m)
	if diags := Analyze(re, env); diags.HasErrors() {
		env = nil
	} else {
		for _, block := range env.Blocks {
			skip[block] = true
		}
		block := RewriteEnvironment(re, env)
		if block != nil {
			blocks = append(blocks, block)
		}
	}

	cs := NewComponentSet(m)
	if diags := Analyze(re, cs); diags.HasErrors() {
		cs = nil
	} else {
		for _, block := range cs.Blocks {
			skip[block] = true
		}
		block := RewriteComponents(re, cs)
		if block != nil {
			blocks = append(blocks, block)
		}
	}

	for _, block := range m.Content.Blocks {
		if skip[block] {
			continue
		}
		blocks = append(blocks, hclgen.BlockFromStructure(block))
	}

	body := &hclgen.Body{
		Attributes: hclgen.AttributesFromStructure(m.Content.Attributes),
		Blocks:     blocks,
	}
	return &hclgen.File{
		Body: body,
	}
}

func (_ RewriteBase) RewriteEnvironment(re Rewrite, env *Environment) *hclgen.Block {
	if len(env.Blocks) == 0 {
		return nil
	}
	// TODO: Handle richer environments.
	return hclgen.BlockFromStructure(env.Blocks[0])
}

func (_ RewriteBase) RewriteComponents(re Rewrite, cs *ComponentSet) *hclgen.Block {
	blocks := make([]*hclgen.Block, len(cs.Components))
	for i, c := range cs.Components {
		blocks[i] = RewriteComponent(re, c)
	}
	return &hclgen.Block{
		Type: "components",
		Body: &hclgen.Body{
			Blocks: blocks,
		},
	}
}

func (_ RewriteBase) RewriteComponent(re Rewrite, c *Component) *hclgen.Block {
	return hclgen.BlockFromSyntax(c.Source)
}
