package container

import (
	"path"
	"testing"

	"github.com/deref/exo/internal/providers/docker/compose"
	"github.com/docker/docker/api/types/mount"
	"github.com/stretchr/testify/assert"
)

func TestMakeMount(t *testing.T) {
	workspaceRoot := "/home/test/app"
	homeDir := "/home/test"
	makeMount := func(vm compose.VolumeMount) mount.Mount {
		res, err := makeMountFromVolumeMount(workspaceRoot, homeDir, vm)
		assert.NoError(t, err)
		return res
	}

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeVolume,
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMount{
		VolumeMountLongForm: compose.VolumeMountLongForm{
			Type:   "volume",
			Target: "/home/node/app",
		},
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: workspaceRoot + "/testing",
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMount{
		VolumeMountLongForm: compose.VolumeMountLongForm{
			Type:   "bind",
			Source: "./testing",
			Target: "/home/node/app",
		},
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/testing",
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMount{
		VolumeMountLongForm: compose.VolumeMountLongForm{
			Type:   "bind",
			Source: "/testing",
			Target: "/home/node/app",
		},
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: path.Join(homeDir, "testing"),
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMount{
		VolumeMountLongForm: compose.VolumeMountLongForm{
			Type:   "bind",
			Source: "~/testing",
			Target: "/home/node/app",
		},
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeVolume,
		Source: "testing",
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMount{
		VolumeMountLongForm: compose.VolumeMountLongForm{
			Type:   "volume",
			Source: "testing",
			Target: "/home/node/app",
		},
	}))
}
