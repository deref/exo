package container

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/docker/docker/api/types/mount"
)

func makeMountFromVolumeMount(workspaceRoot, userHomeDir string, va compose.VolumeMount) (mount.Mount, error) {
	var mountType mount.Type
	var bindOptions *mount.BindOptions
	var volumeOptions *mount.VolumeOptions
	var tmpfsOptions *mount.TmpfsOptions

	switch va.Type {
	case "bind":
		mountType = mount.TypeBind
		bindOptions = &mount.BindOptions{
			Propagation:  mount.Propagation(va.Bind.Propagation),
			NonRecursive: !va.Bind.CreateHostPath,
		}

	case "volume":
		mountType = mount.TypeVolume
		volumeOptions = &mount.VolumeOptions{
			NoCopy: va.Volume.Nocopy,
		}

	case "tmpfs":
		mountType = mount.TypeTmpfs
		tmpfsOptions = &mount.TmpfsOptions{
			SizeBytes: va.Tmpfs.Size,
		}

	default:
		return mount.Mount{}, fmt.Errorf("unsupported mount type: %q", va.Type)
	}

	source := va.Source
	if strings.HasPrefix(source, ".") {
		var err error
		source, err = filepath.Abs(filepath.Join(workspaceRoot, source))
		if err != nil {
			return mount.Mount{}, err
		}
	} else if strings.HasPrefix(source, "~/") {
		source = filepath.Join(userHomeDir, source[2:])
	}

	return mount.Mount{
		Type:          mountType,
		Source:        source,
		Target:        va.Target,
		ReadOnly:      va.ReadOnly,
		BindOptions:   bindOptions,
		VolumeOptions: volumeOptions,
		TmpfsOptions:  tmpfsOptions,
	}, nil
}
