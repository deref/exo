package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/goccy/go-yaml"
)

var Version = "0.1"

type Manifest struct {
	Exo        string      `json:"exo"`
	Components []Component `json:"components"`
}

type ComponentSpec string

func (spec ComponentSpec) MarshalYAML() (interface{}, error) {
	return spec, nil
}

func (spec ComponentSpec) MarshalJSON() ([]byte, error) {
	var d interface{}
	if err := yaml.Unmarshal([]byte(spec), &d); err != nil {
		return nil, fmt.Errorf("spec is not valid yaml")
	}

	bs, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("marshalling component spec json: %w", err)
	}
	return bs, nil
}

func (spec *ComponentSpec) UnmarshalJSON(b []byte) error {
	s := string(b)
	if !yamlutil.IsValid(s) {
		return fmt.Errorf("component spec is not valid yaml")
	}

	*spec = ComponentSpec(s)
	return nil
}

type Component struct {
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Spec      ComponentSpec `json:"spec"` // TODO: Custom unmarshalling to allow convenient json representation.
	DependsOn []string      `json:"depends_on"`
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

func (lr LoadResult) AddUnsupportedFeatureWarning(featureName, explanation string) LoadResult {
	warning := fmt.Sprintf("unsupported feature %s: %s", featureName, explanation)
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
	manifest := Manifest{}
	if err := json.Unmarshal(bs, &manifest); err != nil {
		return LoadResult{Err: err}
	}
	return LoadResult{Manifest: &manifest}
}

func Generate(w io.Writer, manifest *Manifest) error {
	bytes, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("outputting manifest file: %w", err)
	}
	_, err = w.Write(bytes)
	return err
}
