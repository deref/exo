package compose

import "testing"

func TestServiceDependencyYAML(t *testing.T) {
	testYAML(t, "short", `service`, ServiceDependency{
		IsShortSyntax: true,
		Service:       "service",
		ServiceDependencyLongForm: ServiceDependencyLongForm{
			Condition: "service_started",
		},
	})
	testYAML(t, "long", `
condition: service_healthy
`, ServiceDependency{
		ServiceDependencyLongForm: ServiceDependencyLongForm{
			Condition: "service_healthy",
		},
	})
	testYAML(t, "seq", `
- foo
- bar
`, ServiceDependencies{
		Style: SeqStyle,
		Items: []ServiceDependency{
			ServiceDependency{
				IsShortSyntax: true,
				Service:       "foo",
				ServiceDependencyLongForm: ServiceDependencyLongForm{
					Condition: "service_started",
				},
			},
			ServiceDependency{
				IsShortSyntax: true,
				Service:       "bar",
				ServiceDependencyLongForm: ServiceDependencyLongForm{
					Condition: "service_started",
				},
			},
		},
	})
	testYAML(t, "map", `
foo: {}
bar:
  condition: service_healthy
`, ServiceDependencies{
		Style: MapStyle,
		Items: []ServiceDependency{
			ServiceDependency{
				Service: "foo",
			},
			ServiceDependency{
				Service: "bar",
				ServiceDependencyLongForm: ServiceDependencyLongForm{
					Condition: "service_healthy",
				},
			},
		},
	})
}
