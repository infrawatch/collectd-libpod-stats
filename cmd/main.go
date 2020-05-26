package main

import (
	"github.com/collectd/go-collectd/plugin"
	"github.com/pleimer/collectd-libpod-stats/pkg/stats"
)

func init() {
	plugin.RegisterRead("podmanstats", stats.PodmanStats{})
	//plugin.RegisterConfig("Service", stats.Service{})
}

func main() {}
