package main

import (
	"collectd.org/plugin"
	stats "github.com/pleimer/collectd-podman-stats/pkg"
)

func init() {
	plugin.RegisterRead("podman-stats", stats.PodmanStats{})
}

func main() {}
