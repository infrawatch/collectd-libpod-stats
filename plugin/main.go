package main

import (
	"collectd.org/plugin"
)

func init() {
	plugin.RegisterRead("libpodstats", NewLibpodStats())
}

func main() {}
