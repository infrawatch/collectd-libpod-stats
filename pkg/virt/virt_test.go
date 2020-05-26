package virt

import (
	"testing"

	"github.com/pleimer/collectd-libpod-stats/pkg/assert"
)

func TestRunningContainers(t *testing.T) {
	containers, err := GetContainers()
	assert.Ok(t, err)

	for key := range containers {
		t.Log(key)
	}
}

func TestContainerStats(t *testing.T) {
	containers, err := GetContainers()
	assert.Ok(t, err)

	for key, val := range containers {
		stats, err := ContainerStats(key)

		if err != nil {
			continue
		}
		t.Logf("Container: %s %s", val.Names, string(stats))
	}
}
