package compose

import (
	"testing"

	"github.com/deref/exo/internal/util/yamlutil"
	"github.com/stretchr/testify/assert"
)

func TestPortSyntax(t *testing.T) {
	assertParsed := func(expected PortMapping, short string) {
		actual, err := ParsePortMapping(short)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	assertParsed(PortMapping{
		Target: "3000",
	}, "3000")
	assertParsed(PortMapping{
		Target: "3000-3005",
	}, "3000-3005")
	assertParsed(PortMapping{
		Published: "9090-9091",
		Target:    "8080-8081",
	}, "9090-9091:8080-8081")
	assertParsed(PortMapping{
		Published: "49100",
		Target:    "22",
	}, "49100:22")
	assertParsed(PortMapping{
		HostIP:    "127.0.0.1",
		Published: "8001",
		Target:    "8001",
	}, "127.0.0.1:8001:8001")
	assertParsed(PortMapping{
		HostIP:    "127.0.0.1",
		Published: "5000-5010",
		Target:    "5000-5010",
	}, "127.0.0.1:5000-5010:5000-5010")
	assertParsed(PortMapping{
		Published: "6060",
		Target:    "6060",
		Protocol:  "udp",
	}, "6060:6060/udp")
	assertParsed(PortMapping{
		HostIP:    "127.0.0.1",
		Published: "",
		Target:    "65432",
		Protocol:  "udp",
	}, "127.0.0.1::65432/udp")
	assertParsed(PortMapping{
		HostIP:    "2a01:4f8:440c:5e94::1",
		Published: "53",
		Target:    "53",
		Protocol:  "udp",
	}, "2a01:4f8:440c:5e94::1:53:53/udp")
	assertParsed(PortMapping{
		HostIP:    "::",
		Published: "53",
		Target:    "53",
		Protocol:  "udp",
	}, ":::53:53/udp")
}

func TestPortYaml(t *testing.T) {
	type Data struct {
		Ports PortMappings `yaml:"ports"`
	}
	if false { // XXX
		var actual Data
		yamlutil.MustUnmarshalString(`
ports: '1:2,3.4.5.6:7:8/udp',
`, &actual)
		assert.Equal(t, PortMappings{
			PortMapping{
				Published: "1",
				Target:    "2",
			},
			PortMapping{
				HostIP:    "3.4.5.6",
				Published: "7",
				Target:    "8",
				Protocol:  "udp",
			},
		}, actual.Ports)
	}
	{
		var actual Data
		yamlutil.MustUnmarshalString(`
ports:
  - '1:2'
  - '3.4.5.6:7:8/udp'
  - target: 80
    host_ip: 127.0.0.1
    published: 8080
    protocol: tcp
    mode: host
`, &actual)
		assert.Equal(t, PortMappings{
			PortMapping{
				Published: "1",
				Target:    "2",
			},
			PortMapping{
				HostIP:    "3.4.5.6",
				Published: "7",
				Target:    "8",
				Protocol:  "udp",
			},
			PortMapping{
				HostIP:    "127.0.0.1",
				Published: "8080",
				Target:    "80",
				Protocol:  "tcp",
				Mode:      "host",
			},
		}, actual.Ports)
	}
}
