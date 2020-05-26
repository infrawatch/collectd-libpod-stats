package stats

import (
	"context"
	"fmt"
	"time"

	"collectd.org/api"
	"github.com/collectd/go-collectd/plugin"
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

// PodmanStats gather container stats from podman
type PodmanStats struct{}

func (PodmanStats) Read(ctx context.Context) error {
	vl := &api.ValueList{
		Identifier: api.Identifier{
			Host:   "localhost",
			Plugin: "podmanstats",
			Type:   "gauge",
		},
		Time:     time.Now(),
		Interval: 10 * time.Second,
		Values:   []api.Value{api.Gauge(42)},
		DSNames:  []string{"value"},
	}

	plugin.Info("I REALLY tried to execute this plugin")

	if err := plugin.Write(ctx, vl); err != nil {
		return fmt.Errorf("plugin.Write: %w", err)
	}

	return nil
}
