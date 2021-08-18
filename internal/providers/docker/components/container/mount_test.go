package container

import (
	"os/user"
	"path"
	"testing"

	"github.com/docker/docker/api/types/mount"
	"github.com/stretchr/testify/assert"
)

func TestMakeMount(t *testing.T) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := user.HomeDir

	workspaceRoot := "/home/test/app"
	makeMount := func(volumeString string) mount.Mount {
		res, err := makeMountFromVolumeString(workspaceRoot, volumeString)
		assert.NoError(t, err)
		return res
	}

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/home/node/app",
	}, makeMount("/home/node/app"))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: workspaceRoot + "/testing",
		Target: "/home/node/app",
	}, makeMount("./testing:/home/node/app"))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: "/testing",
		Target: "/home/node/app",
	}, makeMount("/testing:/home/node/app"))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeBind,
		Source: path.Join(homeDir, "testing"),
		Target: "/home/node/app",
	}, makeMount("~/testing:/home/node/app"))

	assert.Equal(t, mount.Mount{
		Type:   mount.TypeVolume,
		Source: "testing",
		Target: "/home/node/app",
	}, makeMount("testing:/home/node/app"))
}
