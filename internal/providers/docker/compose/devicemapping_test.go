package compose

import "testing"

func TestDeviceMapping(t *testing.T) {
	testYAML(t, "no_perms", `/dev/ttyUSB0:/dev/ttyUSB1`, DeviceMapping{
		PathOnHost:      "/dev/ttyUSB0",
		PathInContainer: "/dev/ttyUSB1",
	})
	testYAML(t, "with_perms", `/dev/sda:/dev/xvda:rwm`, DeviceMapping{
		PathOnHost:        "/dev/sda",
		PathInContainer:   "/dev/xvda",
		CgroupPermissions: "rwm",
	})
}
