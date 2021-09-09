package compose

import (
	"errors"
	"strings"
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

func (dm *DeviceMapping) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var dmStr string
	if err := unmarshal(&dmStr); err != nil {
		return err
	}

	segments := strings.Split(dmStr, ":")
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
