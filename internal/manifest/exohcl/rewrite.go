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
	RewriteSecrets(Rewrite, *Secrets) *hclgen.Block
	RewriteComponents(Rewrite, *ComponentSet) *hclgen.Block
	RewriteComponent(Rewrite, *Component) *hclgen.Block
}

func RewriteManifest(re Rewrite, m *Manifest) *hclgen.File {
	return re.RewriteManifest(re, m)
}

func RewriteEnvironment(re Rewrite, env *Environment) *hclgen.Block {
	return re.RewriteEnvironment(re, env)
}

func RewriteSecrets(re Rewrite, s *Secrets) *hclgen.Block {
	return re.RewriteSecrets(re, s)
}

func RewriteComponents(re Rewrite, cs *ComponentSet) *hclgen.Block {
	return re.RewriteComponents(re, cs)
}

func RewriteComponent(re Rewrite, c *Component) *hclgen.Block {
	return re.RewriteComponent(re, c)
}

// RewriteBase is intended embedded in Rewrite implementations to provide
// default method implementations. The default implementations rewrite
// structures in to normal form.
type RewriteBase struct{}

// Normalize lifts RewriteBase to fully satisfy the Rewrite interface.
type Normalize struct {
	context.Context
	RewriteBase
}

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
	children := make([]*hclgen.Block, 0, len(env.Secrets))
	for _, secretsIn := range env.Secrets {
		secretsOut := RewriteSecrets(re, secretsIn)
		if secretsOut == nil {
			continue
		}
		children = append(children, secretsOut)
	}
	if len(env.Attributes) == 0 && len(children) == 0 {
		return nil
	}
	return &hclgen.Block{
		Type:   "environment",
		Labels: env.Blocks[0].Labels, // XXX OK to discard labels from extra blocks?
		Body: &hclgen.Body{
			Attributes: env.Attributes,
			Blocks:     children,
		},
	}
}

func (_ RewriteBase) RewriteSecrets(re Rewrite, s *Secrets) *hclgen.Block {
	return &hclgen.Block{
		Type:   "secrets",
		Labels: s.Block.Labels,
		Body:   hclgen.BodyFromStructure(s.Block.Body),
	}
}

func (_ RewriteBase) RewriteComponents(re Rewrite, cs *ComponentSet) *hclgen.Block {
	n := len(cs.Components)
	if n == 0 {
		return nil
	}
	blocks := make([]*hclgen.Block, n)
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
