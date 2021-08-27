package manifest

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var Version = "0.1"

type Manifest struct {
	Exo        string      `hcl:"exo"`
	Components []Component `hcl:"component,block"`
}

type Component struct {
	Name      string   `hcl:"name,label"`
	Type      string   `hcl:"type,label"`
	Spec      string   `hcl:"spec"` // TODO: Custom unmarshalling to allow convenient json representation.
	DependsOn []string `hcl:"depends_on"`
}

func NewManifest() *Manifest {
	return &Manifest{
		Exo: Version,
	}
}

type LoadResult struct {
	Manifest *Manifest
	Warnings []string
	Err      error
}

func (lr LoadResult) AddRenameWarning(originalName, newName string) LoadResult {
	warning := fmt.Sprintf("invalid name: %q, renamed to: %q", originalName, newName)
	lr.Warnings = append(lr.Warnings, warning)
	return lr
}

type loader struct{}

var Loader = loader{}

func (l loader) Load(r io.Reader) LoadResult {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return LoadResult{Err: err}
	}
	return ReadBytes(bs)
}

func ReadBytes(bs []byte) LoadResult {
	var manifest Manifest
	evalCtx := &hcl.EvalContext{
		Functions: map[string]function.Function{
			"jsonencode": stdlib.JSONEncodeFunc,
		},
	}
	if err := hclsimple.Decode("exo.hcl", bs, evalCtx, &manifest); err != nil {
		return LoadResult{Err: err}
	}
	return LoadResult{Manifest: &manifest}
}

func Generate(w io.Writer, manifest *Manifest) error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(manifest, f.Body())
	_, err := f.WriteTo(w)
	return err
}
