package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePortMappings(t *testing.T) {
	assertParsed := func(expected PortMappingLongForm, short string) {
		actual, err := ParsePortMapping(short)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	assertParsed(PortMappingLongForm{
		Target: PortRange{
			Min: 3000,
			Max: 3000,
		},
	}, "3000")
	assertParsed(PortMappingLongForm{
		Target: PortRange{
			Min: 3000,
			Max: 3005,
		},
	}, "3000-3005")
	assertParsed(PortMappingLongForm{
		Published: PortRange{
			Min: 9090,
			Max: 9091,
		},
		Target: PortRange{
			Min: 8080,
			Max: 8081,
		},
	}, "9090-9091:8080-8081")
	assertParsed(PortMappingLongForm{
		Published: PortRange{
			Min: 49100,
			Max: 49100,
		},
		Target: PortRange{
			Min: 22,
			Max: 22,
		},
	}, "49100:22")
	assertParsed(PortMappingLongForm{
		HostIP: "127.0.0.1",
		Published: PortRange{
			Min: 8001,
			Max: 8001,
		},
		Target: PortRange{
			Min: 8001,
			Max: 8001,
		},
	}, "127.0.0.1:8001:8001")
	assertParsed(PortMappingLongForm{
		HostIP: "127.0.0.1",
		Published: PortRange{
			Min: 5000,
			Max: 5010,
		},
		Target: PortRange{
			Min: 5000,
			Max: 5010,
		},
	}, "127.0.0.1:5000-5010:5000-5010")
	assertParsed(PortMappingLongForm{
		Published: PortRange{
			Min: 6060,
			Max: 6060,
		},
		Target: PortRange{
			Min: 6060,
			Max: 6060,
		},
		Protocol: "udp",
	}, "6060:6060/udp")
	assertParsed(PortMappingLongForm{
		HostIP:    "127.0.0.1",
		Published: PortRange{},
		Target: PortRange{
			Min: 65432,
			Max: 65432,
		},
		Protocol: "udp",
	}, "127.0.0.1::65432/udp")
	assertParsed(PortMappingLongForm{
		HostIP: "2a01:4f8:440c:5e94::1",
		Published: PortRange{
			Min: 53,
			Max: 53,
		},
		Target: PortRange{
			Min: 53,
			Max: 53,
		},
		Protocol: "udp",
	}, "2a01:4f8:440c:5e94::1:53:53/udp")
	assertParsed(PortMappingLongForm{
		HostIP: "::",
		Published: PortRange{
			Min: 53,
			Max: 53,
		},
		Target: PortRange{
			Min: 53,
			Max: 53,
		},
		Protocol: "udp",
	}, ":::53:53/udp")
}

func TestPortYAML(t *testing.T) {
	testYAML(t, "short_target", `4000`, PortMapping{
		IsShortForm: true,
		PortMappingLongForm: PortMappingLongForm{
			Target: PortRange{
				Min: 4000,
				Max: 4000,
			},
		},
	})
	testYAML(t, "short_target_range", `4000-4001`, PortMapping{
		IsShortForm: true,
		PortMappingLongForm: PortMappingLongForm{
			Target: PortRange{
				Min: 4000,
				Max: 4001,
			},
		},
	})
	testYAML(t, "long", `
target: 4000-4001
published: 5000-5001
host_ip: 1.2.3.4
protocol: tcp
mode: host
`, PortMapping{
		PortMappingLongForm: PortMappingLongForm{
			Target: PortRange{
				Min: 4000,
				Max: 4001,
			},
			Published: PortRange{
				Min: 5000,
				Max: 5001,
			},
			HostIP:   "1.2.3.4",
			Protocol: "tcp",
			Mode:     "host",
		},
	})
	testYAML(t, "multiple", `
- 3000
- published: 4000
`, PortMappings{
		PortMapping{
			IsShortForm: true,
			PortMappingLongForm: PortMappingLongForm{
				Target: PortRange{
					Min: 3000,
					Max: 3000,
				},
			},
		},
		PortMapping{
			PortMappingLongForm: PortMappingLongForm{
				Published: PortRange{
					Min: 4000,
					Max: 4000,
				},
			},
		},
	})
	testYAML(t, "range_with_protocol", `1000-2000/tcp`, PortRangeWithProtocol{
		PortRange: PortRange{
			Min: 1000,
			Max: 2000,
		},
		Protocol: "tcp",
	})
	testYAML(t, "range_with_blank_protocol", `1000-2000`, PortRangeWithProtocol{
		PortRange: PortRange{
			Min: 1000,
			Max: 2000,
		},
	})
}
