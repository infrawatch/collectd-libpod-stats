package main

import (
	"context"
	"fmt"
	"time"

	"collectd.org/api"
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
	statMatrix, err := virt.ContainersStats(cgroups.CPUAcctT)
	if err != nil {
		return err
	}
	vl := &api.ValueList{
		Identifier: api.Identifier{
			Host:   "localhost",
			Plugin: "libpodstats",
			Type:   "gauge",
		},
		Time:     time.Now(),
		Interval: 10 * time.Second,
		Values:   []api.Value{api.Counter(statMatrix["qdr"][cgroups.CPUAcctT])},
		DSNames:  []string{"value"},
	}

	plugin.Info("I REALLY tried to execute this plugin")

	if err := plugin.Write(ctx, vl); err != nil {
		return fmt.Errorf("plugin.Write: %w", err)
	}

	return nil
}
