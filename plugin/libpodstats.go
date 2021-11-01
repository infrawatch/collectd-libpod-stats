package main

import (
	"context"
	"fmt"
	"time"

	"collectd.org/api"
	"collectd.org/plugin"
	"github.com/infrawatch/collectd-libpod-stats/pkg/cgroups"
	"github.com/infrawatch/collectd-libpod-stats/pkg/virt"
)

// LibpodStats gathers container resource usage stats from libpod cgroups
type LibpodStats struct {
	handlers map[cgroups.ControlType]handler
}

//NewLibpodStats initialize new libpodstats plugins with handlers
//TODO: generate handlers from plugin config
func NewLibpodStats() *LibpodStats {
	handlers := map[cgroups.ControlType]handler{}

	handlers[cgroups.CPUAcctT] = &cpuHandler{}
	handlers[cgroups.MemoryT] = &memoryHandler{}
	return &LibpodStats{
		handlers: handlers,
	}
}

func (ls *LibpodStats) Read(ctx context.Context) error {
	statMatrix, err := virt.ContainersStats(cgroups.CPUAcctT, cgroups.MemoryT)
	if err != nil {
		return err
	}

	for cLabel, metric := range statMatrix {
		for controlType, stat := range metric {
			vl := &api.ValueList{
				Identifier: api.Identifier{
					Plugin:         "libpodstats",
					PluginInstance: cLabel,
				},
				Time:     time.Now(),
				Interval: 10 * time.Second,
			}

			if _, found := ls.handlers[controlType]; !found {
				return fmt.Errorf("unhandled cgroup type: %s", controlType.String())
			}
			ls.handlers[controlType].populateValueList(stat, vl)

			if err := plugin.Write(ctx, vl); err != nil {
				return fmt.Errorf("plugin.Write: %w", err)
			}
		}
	}
	return nil
}

func init() {
	plugin.RegisterRead("libpodstats", NewLibpodStats())
}

func main() {}
