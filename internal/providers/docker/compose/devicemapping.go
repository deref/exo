package compose

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

type DeviceMapping struct {
	PathOnHost        string
	PathInContainer   string
	CgroupPermissions string
}

func (dm DeviceMapping) MarshalYAML() (interface{}, error) {
	var out strings.Builder
	out.WriteString(dm.PathOnHost)
	out.WriteByte(':')
	out.WriteString(dm.PathInContainer)
	if dm.CgroupPermissions != "" {
		out.WriteByte(':')
		out.WriteString(dm.CgroupPermissions)
	}

	return out.String(), nil
}

func (dm *DeviceMapping) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	segments := strings.Split(s, ":")
	if len(segments) < 2 || len(segments) > 3 {
		return errors.New("device mapping should be in the form: HOST_PATH:CONTAINER_PATH[:CGROUP_PERMISSIONS]")
	}

	dm.PathOnHost = segments[0]
	dm.PathInContainer = segments[1]
	if len(segments) == 3 {
		dm.CgroupPermissions = segments[2]
	}

	return nil
}
