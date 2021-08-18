package container

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/mount"
)

func makeMountFromVolumeString(workspaceRoot, volume string) (mount.Mount, error) {
	// This obviously doesn't handle colons in the path well but it would appear
	// that docker-compose doesn't handle those either.
	volumeParts := strings.Split(volume, ":")
	if len(volumeParts) > 2 {
		return mount.Mount{}, fmt.Errorf("invalid volume string %s", volume)
	}
	if len(volumeParts) == 1 {
		return mount.Mount{
			Type:   mount.TypeBind,
			Source: volumeParts[0],
		}, nil
	}

	source, target := volumeParts[0], volumeParts[1]
	mountType := mount.TypeBind
	if strings.HasPrefix(source, "./") {
		var err error
		source, err = filepath.Abs(filepath.Join(workspaceRoot, source))
		if err != nil {
			return mount.Mount{}, err
		}
	} else if strings.HasPrefix(source, "~/") {
		user, err := user.Current()
		if err != nil {
			return mount.Mount{}, err
		}
		source = filepath.Join(user.HomeDir, source[2:])
	} else if !strings.HasPrefix(source, "/") {
		mountType = mount.TypeVolume
	}

	result := mount.Mount{
		Type:   mountType,
		Source: source,
		Target: target,
	}

	return result, nil
}
