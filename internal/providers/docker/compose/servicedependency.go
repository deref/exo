package compose

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

// ServiceDependenciesTemplate represents either the short (list) form of
// `depends_on` or the full form that indicates the condition under which this
// service can start. Since there are two mutally exclusive representations of
// this structure, custom marshaling and unmarshing is implemented.
type ServiceDependenciesTemplate []ServiceDependencyTemplate

type ServiceDependencies []ServiceDependency

type ServiceDependencyTemplate struct {
	Service   string
	Condition string
}

type ServiceDependency struct {
	Service   string
	Condition string // TODO: validation.
}

func (sd ServiceDependenciesTemplate) MarshalYAML() (interface{}, error) {
	services := make(map[string]interface{}, len(sd))
	for _, service := range sd {
		services[service.Service] = map[string]interface{}{
			"condition": service.Condition,
		}
	}
	return services, nil
}

func (sd *ServiceDependenciesTemplate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asStrings []string
	err := unmarshal(&asStrings)
	if err == nil {
		*sd = make([]ServiceDependencyTemplate, len(asStrings))
		for i, service := range asStrings {
			(*sd)[i] = ServiceDependencyTemplate{
				Service:   service,
				Condition: "service_started",
			}
		}
		return nil
	}

	var slice yaml.MapSlice
	if err := unmarshal(&slice); err != nil {
		return err
	}
	var m map[string]struct {
		Condition string `yaml:"condition,omitempty"`
	}
	if err := unmarshal(&m); err != nil {
		return err
	}

	*sd = make([]ServiceDependencyTemplate, len(slice))
	for i, item := range slice {
		name, ok := item.Key.(string)
		if !ok {
			return fmt.Errorf("expected key to be string, got %T", item.Key)
		}
		dependency := m[name]
		(*sd)[i] = ServiceDependencyTemplate{
			Service:   name,
			Condition: dependency.Condition,
		}
	}

	return nil
}

func (sd *ServiceDependencies) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dependencies ServiceDependenciesTemplate
	if err := unmarshal(&dependencies); err != nil {
		return err
	}
	*sd = make([]ServiceDependency, len(dependencies))
	for i, dependency := range dependencies {
		switch dependency.Condition {
		case "":
			dependency.Condition = "service_started"
		case "service_started", "service_healthy", "service_completed_successfully":
			// Ok.
		default:
			return fmt.Errorf("invalid condition %q for service dependency %q", dependency.Condition, dependency.Service)
		}
		(*sd)[i] = ServiceDependency{
			Service:   dependency.Service,
			Condition: dependency.Condition,
		}
	}
	return nil
}
