package compose

import "testing"

func TestServiceDependencyYAML(t *testing.T) {
	testYAML(t, "short", `service`, ServiceDependency{
		IsShortSyntax: true,
		Service:       MakeString("service"),
	})
	testYAML(t, "long", `
condition: service_healthy
`, ServiceDependency{
		ServiceDependencyLongForm: ServiceDependencyLongForm{
			Condition: MakeString("service_healthy"),
		},
	})
	testYAML(t, "seq", `
- foo
- bar
`, ServiceDependencies{
		Style: SeqStyle,
		Items: []ServiceDependency{
			{
				IsShortSyntax:             true,
				Service:                   MakeString("foo"),
				ServiceDependencyLongForm: ServiceDependencyLongForm{},
			},
			{
				IsShortSyntax:             true,
				Service:                   MakeString("bar"),
				ServiceDependencyLongForm: ServiceDependencyLongForm{},
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
			{
				Service: MakeString("foo"),
			},
			{
				Service: MakeString("bar"),
				ServiceDependencyLongForm: ServiceDependencyLongForm{
					Condition: MakeString("service_healthy"),
				},
			},
		},
	})
	assertInterpolated(t, map[string]string{"short": "service"}, `
- ${short}
`, ServiceDependencies{
		Style: SeqStyle,
		Items: []ServiceDependency{
			{
				IsShortSyntax: true,
				Service:       MakeString("${short}").WithValue("service"),
			},
		},
	})
	assertInterpolated(t, map[string]string{"condition": "service_healthy"}, `
service:
  condition: ${condition}
`, ServiceDependencies{
		Style: MapStyle,
		Items: []ServiceDependency{
			{
				Service: MakeString("service"),
				ServiceDependencyLongForm: ServiceDependencyLongForm{
					Condition: MakeString("${condition}").WithValue("service_healthy"),
				},
			},
		},
	})
}
