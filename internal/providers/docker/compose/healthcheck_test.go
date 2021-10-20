package compose

import "testing"

func TestHealthcheckYAML(t *testing.T) {
	assertInterpolated(t, map[string]string{"one": "1"}, `
retries: ${one}
`, Healthcheck{
		Retries: Int{
			String: MakeString("${one}").WithValue("1"),
			Value:  1,
		},
	})
}
