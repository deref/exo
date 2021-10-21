package compose

import "testing"

func TestLinkL(t *testing.T) {
	testYAML(t, "unaliased", `service`, Link{
		String:  MakeString("service"),
		Service: "service",
		Alias:   "service",
	})
	testYAML(t, "alias", `service:alias`, Link{
		String:  MakeString("service:alias"),
		Service: "service",
		Alias:   "alias",
	})
	assertInterpolated(t, map[string]string{"service": "SERVICE"}, `${service}:alias`, Link{
		String:  MakeString("${service}:alias").WithValue("SERVICE:alias"),
		Service: "SERVICE",
		Alias:   "alias",
	})
}
