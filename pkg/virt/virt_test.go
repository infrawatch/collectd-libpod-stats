package virt

import (
	"testing"

	"github.com/pleimer/collectd-libpod-stats/pkg/assert"
	"github.com/pleimer/collectd-libpod-stats/pkg/cgroups"
)

type userPaths struct {
	root    string
	nonroot string
}

type cgrps struct {
	version1 userPaths
	version2 userPaths
}

var testPathMatrix cgrps = cgrps{
	version1: userPaths{
		root:    "/sys/fs/cgroup/cpu/machine.slice/libpod-123abc.scope",
		nonroot: "",
	},
	version2: userPaths{
		root:    "/sys/fs/cgroup/machine.slice/libpod-123abc.scope",
		nonroot: "/sys/fs/cgroup/user.slice/user-1000.slice/user@1000.service/user.slice/libpod-123abc.scope",
	},
}

func TestCgroupPath(t *testing.T) {
	t.Run("root cgroup v1", func(t *testing.T) {
		path, err := genContainerCgroupPath(cgroups.CPUAcctT, "123abc", false, 0)
		assert.Ok(t, err)
		assert.Equals(t, testPathMatrix.version1.root, path)
	})

	t.Run("nonroot cgroup v1", func(t *testing.T) {
		_, err := genContainerCgroupPath(cgroups.CPUAcctT, "123abc", false, 1000)
		assert.Assert(t, err != nil, "path generation should fail for cgroups v1 rootless")
	})

	t.Run("root cgroup v2", func(t *testing.T) {
		path, err := genContainerCgroupPath(cgroups.CPUAcctT, "123abc", true, 0)
		assert.Ok(t, err)
		assert.Equals(t, testPathMatrix.version2.root, path)
	})

	t.Run("nonroot cgroup v2", func(t *testing.T) {
		path, err := genContainerCgroupPath(cgroups.CPUAcctT, "123abc", true, 1000)
		assert.Ok(t, err)
		assert.Equals(t, testPathMatrix.version2.nonroot, path)
	})
}
