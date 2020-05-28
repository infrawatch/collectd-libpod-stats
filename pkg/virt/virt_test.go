package virt

import (
	"testing"

	"github.com/pleimer/collectd-libpod-stats/pkg/assert"
	"github.com/pleimer/collectd-libpod-stats/pkg/cgroups"
)

func TestContainerStats(t *testing.T) {
	statMatrix, err := ContainersStats(cgroups.CPUAcctT, cgroups.MemoryT)
	assert.Ok(t, err)
	t.Logf("Qdr: %v", statMatrix["qdr"])
}
