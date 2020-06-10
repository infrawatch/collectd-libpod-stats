package main

import (
	"github.com/collectd/go-collectd/plugin"
)

func init() {
	plugin.RegisterRead("libpodstats", NewLibpodStats())
	//plugin.RegisterConfig("Service", stats.Service{})
}

func main() {}
