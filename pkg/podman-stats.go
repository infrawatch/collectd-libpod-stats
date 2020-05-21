package stats

import (
	"context"
	"fmt"
	"time"

	"collectd.org/api"
	"collectd.org/plugin"
)

// PodmanStats gather container stats from podman
type PodmanStats struct{}

func (PodmanStats) Read(ctx context.Context) error {
	vl := &api.ValueList{
		Identifier: api.Identifier{
			Host:   "example.com",
			Plugin: "goplug",
			Type:   "gauge",
		},
		Time:     time.Now(),
		Interval: 10 * time.Second,
		Values:   []api.Value{api.Gauge(42)},
		DSNames:  []string{"value"},
	}
	if err := plugin.Write(ctx, vl); err != nil {
		return fmt.Errorf("plugin.Write: %w", err)
	}

	return nil
}
