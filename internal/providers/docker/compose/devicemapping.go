package compose

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

type DeviceMapping struct {
	String
	PathOnHost        string
	PathInContainer   string
	CgroupPermissions string
}

func (dm DeviceMapping) MarshalYAML() (interface{}, error) {
	if dm.String.Expression != "" {
		return dm.String.Expression, nil
	}

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
	if err := node.Decode(&dm.String); err != nil {
		return err
	}

	_ = dm.Interpolate(nil)
	return nil
}

func (dm *DeviceMapping) Interpolate(env Environment) error {
	if err := dm.String.Interpolate(env); err != nil {
		return err
	}

	segments := strings.Split(dm.String.Value, ":")
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
