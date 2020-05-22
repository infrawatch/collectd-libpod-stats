package main

import (
	"github.com/collectd/go-collectd/plugin"
	stats "github.com/pleimer/collectd-podman-stats/pkg"
)

func init() {
	plugin.RegisterRead("podmanstats", stats.PodmanStats{})
	plugin.RegisterConfig("Service", stats.Service{})
}

func main() {}
