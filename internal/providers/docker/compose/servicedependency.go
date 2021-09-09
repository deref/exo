package compose

import (
	"fmt"
	"sort"
	"strings"
)

// ServiceDependencies represents either the short (list) form of `depends_on` or the full
// form that indicates the condition under which this service can start. Since there are two
// mutally exclusive representations of this structure, custom marshaling and unmarshing is
// implemented.
type ServiceDependencies struct {
	Services      []ServiceDependency
	IsShortSyntax bool
}

type ServiceDependency struct {
	Service   string
	Condition string
}

func (sd ServiceDependencies) MarshalYAML() (interface{}, error) {
	if sd.IsShortSyntax {
		services := make([]string, len(sd.Services))
		for i, service := range sd.Services {
			services[i] = service.Service
		}
		return services, nil
	}

	services := make(map[string]interface{}, len(sd.Services))
	for _, service := range sd.Services {
		services[service.Service] = map[string]interface{}{
			"condition": service.Condition,
		}
	}
	return services, nil
}

func (sd *ServiceDependencies) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var asStrings []string
	if err := unmarshal(&asStrings); err == nil {
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

	var asMap map[string]interface{}
	if err := unmarshal(&asMap); err != nil {
		return err
	}

	sd.Services = make([]ServiceDependency, 0, len(asMap))
	for service, value := range asMap {
		condition := "service_started"
		if spec, ok := value.(map[string]interface{}); ok {
			if specCondition, ok := spec["condition"]; ok {
				condition = specCondition.(string)
			}
		}

		switch condition {
		case "service_started", "service_healthy", "service_completed_successfully":
			// Ok.
		default:
			return fmt.Errorf("invalid condition %q for service dependency %q", condition, service)
		}
		sd.Services = append(sd.Services, ServiceDependency{
			Service:   service,
			Condition: condition,
		})
	}
	sort.Sort(serviceDependenciesSort{sd.Services})

	return nil
}

type serviceDependenciesSort struct {
	dependencies []ServiceDependency
}

func (iface serviceDependenciesSort) Len() int {
	return len(iface.dependencies)
}

func (iface serviceDependenciesSort) Less(i, j int) bool {
	return strings.Compare(iface.dependencies[i].Service, iface.dependencies[j].Service) < 0
}

func (iface serviceDependenciesSort) Swap(i, j int) {
	tmp := iface.dependencies[i]
	iface.dependencies[i] = iface.dependencies[j]
	iface.dependencies[j] = tmp
}
