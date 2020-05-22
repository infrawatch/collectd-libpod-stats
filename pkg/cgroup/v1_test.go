package cgroup

import (
	"testing"

	"github.com/pleimer/collectd-podman-stats/pkg/assert"
)

func TestGetStat(t *testing.T) {
	cgroupv2, err := IsCgroup2UnifiedMode()
	assert.Ok(t, err)
	t.Logf("CgroupV2? %v\n", cgroupv2)
	v1 := V1{}
	stat, err := v1.GetCPU()
	assert.Ok(t, err)
	t.Log(stat)
}
