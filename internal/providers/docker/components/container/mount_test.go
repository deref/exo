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
	makeMount := func(vm compose.VolumeMountLongForm) mount.Mount {
		res, err := makeMountFromVolumeMount(workspaceRoot, homeDir, compose.VolumeMount{
			VolumeMountLongForm: vm,
		})
		assert.NoError(t, err)
		return res
	}

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeVolume,
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMountLongForm{
		Type:   compose.MakeString("volume"),
		Target: compose.MakeString("/home/node/app"),
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: workspaceRoot + "/testing",
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMountLongForm{
		Type:   compose.MakeString("bind"),
		Source: compose.MakeString("./testing"),
		Target: compose.MakeString("/home/node/app"),
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/testing",
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMountLongForm{
		Type:   compose.MakeString("bind"),
		Source: compose.MakeString("/testing"),
		Target: compose.MakeString("/home/node/app"),
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: path.Join(homeDir, "testing"),
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMountLongForm{
		Type:   compose.MakeString("bind"),
		Source: compose.MakeString("~/testing"),
		Target: compose.MakeString("/home/node/app"),
	}))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeVolume,
		Source: "testing",
		Target: "/home/node/app",
	}, makeMount(compose.VolumeMountLongForm{
		Type:   compose.MakeString("volume"),
		Source: compose.MakeString("testing"),
		Target: compose.MakeString("/home/node/app"),
	}))
}
