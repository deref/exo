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
}
