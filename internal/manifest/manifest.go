package manifest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/deref/exo/internal/util/yamlutil"
	"gopkg.in/yaml.v3"
)

var Version = "0.1"

type Manifest struct {
	original   *yaml.Node
	Exo        string      `json:"exo"`
	Components []Component `json:"components"`
}

// FIXME: should settle on one yaml library. goccy/go-yaml didn't appear to
// expose the necessary Node functionality, specifically the ability to
// deserialise a node. That meant that whilst you could add comments to a YAML
// block before serialisation, it didn't seem possible to preserve comments when
// deserialising.

func (m *Manifest) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type innerManifestType Manifest
	innerManifest := (*innerManifestType)(m)
	if err := unmarshal(innerManifest); err != nil {
		return fmt.Errorf("unmarshalling manifest yaml: %w", err)
	}
	*m = Manifest(*innerManifest)

	m.original = &yaml.Node{}
	if err := unmarshal(m.original); err != nil {
		return fmt.Errorf("unmarshalling raw manifest yaml: %w", err)
	}
	return nil
}

func (spec Manifest) MarshalYAML() ([]byte, error) {
	type innerManifest Manifest
	bs, err := yaml.Marshal(innerManifest(spec))
	if err != nil {
		return nil, fmt.Errorf("marshalling manifest yaml: %w", err)
	}

	node := &yaml.Node{}
	if err := yaml.Unmarshal(bs, node); err != nil {
		return nil, fmt.Errorf("unmarshalling marshalled manifest yaml: %w", err)
	}
	yamlutil.MergeFormatting(spec.original, node)

	result, err := yaml.Marshal(node)
	if err != nil {
		return nil, fmt.Errorf("marshalling node: %w", err)
	}
	return result, nil
}

type ComponentSpec string

func (spec ComponentSpec) MarshalYAML() (interface{}, error) {
	var d interface{}
	if err := yaml.Unmarshal([]byte(spec), &d); err != nil {
		return nil, fmt.Errorf("spec is not valid yaml: %w", err)
	}

	return d, nil
}

func (spec *ComponentSpec) UnmarshalYAML(value *yaml.Node) error {
	bs, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshalling component spec for unmarshalling: %w", err)
	}

	s := string(bs)
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
	DependsOn []string      `json:"depends_on,omitempty"`
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
	if err := yaml.Unmarshal(bs, &manifest); err != nil {
		return LoadResult{Err: fmt.Errorf("unmarshalling manifest yaml: %w", err)}
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
