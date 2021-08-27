package compose

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

// ServiceDependencies represents either the short (list) form of `depends_on` or the full
// form that indicates the condition under which this service can start. Since there are two
// mutally exclusive representations of this structure, marshaling and unmarshing must be done
// manually, and no yaml struct tags are used.
type ServiceDependencies struct {
	Services      []ServiceDependency
	IsShortSyntax bool
}

type ServiceDependency struct {
	Service   string
	Condition string
}

func (sd *ServiceDependencies) UnmarshalYAML(b []byte) error {
	var asStrings []string
	if err := yaml.Unmarshal(b, &asStrings); err == nil {
		sd.IsShortSyntax = true
		sd.Services = make([]ServiceDependency, len(asStrings))
		for i, service := range asStrings {
			sd.Services[i] = ServiceDependency{
				Service:   service,
				Condition: "service_started",
			}
		}
		return nil
	}

	asMap := make(map[string]struct {
		Condition string `yaml:"condition"`
	})
	if err := yaml.Unmarshal(b, &asMap); err != nil {
		return err
	}

	sd.Services = make([]ServiceDependency, 0, len(asMap))
	for service, spec := range asMap {
		switch spec.Condition {
		case "service_started", "service_healthy", "service_completed_successfully":
			// Ok.
		case "":
			spec.Condition = "service_started"
		default:
			return fmt.Errorf("invalid condition %q for service dependency %q", spec.Condition, service)
		}
		sd.Services = append(sd.Services, ServiceDependency{
			Service:   service,
			Condition: spec.Condition,
		})
	}

	return nil
}
