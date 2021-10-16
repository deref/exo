package compose

import "testing"

func TestDeviceMapping(t *testing.T) {
	testYAML(t, "no_perms", `/dev/ttyUSB0:/dev/ttyUSB1`, DeviceMapping{
		String:          MakeString(`/dev/ttyUSB0:/dev/ttyUSB1`),
		PathOnHost:      "/dev/ttyUSB0",
		PathInContainer: "/dev/ttyUSB1",
	})
	testYAML(t, "with_perms", `/dev/sda:/dev/xvda:rwm`, DeviceMapping{
		String:            MakeString(`/dev/sda:/dev/xvda:rwm`),
		PathOnHost:        "/dev/sda",
		PathInContainer:   "/dev/xvda",
		CgroupPermissions: "rwm",
	})
	assertInterpolated(t, map[string]string{"whole": "/dev/sda:/dev/xvda:rwm"}, `${whole}`, DeviceMapping{
		String:            MakeString(`${whole}`).WithValue(`/dev/sda:/dev/xvda:rwm`),
		PathOnHost:        "/dev/sda",
		PathInContainer:   "/dev/xvda",
		CgroupPermissions: "rwm",
	})
	assertInterpolated(t, map[string]string{"part": "/dev/sda"}, `${part}:/dev/xvda:rwm`, DeviceMapping{
		String:            MakeString(`${part}:/dev/xvda:rwm`).WithValue(`/dev/sda:/dev/xvda:rwm`),
		PathOnHost:        "/dev/sda",
		PathInContainer:   "/dev/xvda",
		CgroupPermissions: "rwm",
	})
}
