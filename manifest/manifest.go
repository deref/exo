package manifest

import (
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
	Name string `hcl:"name,label"`
	Type string `hcl:"type,label"`
	Spec string `hcl:"spec"` // TODO: Custom unmarshalling to allow convenient json representation.
}

func NewManifest() *Manifest {
	return &Manifest{
		Exo: Version,
	}
}

func Read(r io.Reader) (*Manifest, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ReadBytes(bs)
}

func ReadBytes(bs []byte) (*Manifest, error) {
	var manifest Manifest
	evalCtx := &hcl.EvalContext{
		Functions: map[string]function.Function{
			"jsonencode": stdlib.JSONEncodeFunc,
		},
	}
	if err := hclsimple.Decode("exo.hcl", bs, evalCtx, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func Generate(w io.Writer, manifest *Manifest) error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(manifest, f.Body())
	_, err := f.WriteTo(w)
	return err
}
