package main

import (
	"context"
	"fmt"
	"time"

	"collectd.org/api"
	"collectd.org/exec"
	"github.com/collectd/go-collectd/plugin"
	"github.com/pleimer/collectd-libpod-stats/pkg/cgroups"
	"github.com/pleimer/collectd-libpod-stats/pkg/virt"
)

type Service struct {
}

// func (Service) Configure(ctx context.Context, block config.Block) error {
// 	configMap := make(map[string]interface{})
// 	err := block.Unmarshal(&configMap)
// 	if err != nil {
// 		return err
// 	}
// 	for key, val := range configMap {
// 		fmt.Printf("%s:%s\n", key, val)
// 	}
// 	return nil
// }

// LibpodStats gather container stats from podman
type LibpodStats struct{}

func (LibpodStats) Read(ctx context.Context) error {
	statMatrix, err := virt.ContainersStats(cgroups.CPUAcctT, cgroups.MemoryT)
	if err != nil {
		return err
	}

	for cName, metric := range statMatrix {
		var vl *api.ValueList
		for controller, stat := range metric {
			fmt.Printf("Controller: %s, stat: %d\n", controller, stat)
			vl = &api.ValueList{
				Identifier: api.Identifier{
					Host:           exec.Hostname(),
					Plugin:         "libpodstats",
					PluginInstance: cName,
					Type:           controller.String(),
				},
				Time:     time.Now(),
				Interval: 10 * time.Second,
				Values:   []api.Value{stat},
			}

			if err := plugin.Write(ctx, vl); err != nil {
				return fmt.Errorf("plugin.Write: %w", err)
			}
		}
	}
	return nil
}
