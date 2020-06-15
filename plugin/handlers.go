package main

import (
	"time"

	"collectd.org/api"
)

type handler interface {
	populateValueList(uint64, *api.ValueList)
}

type pair struct {
	systemTime uint64
	cpuTime    uint64
}

//cpuHandler stat gathering and formatting for cpu cgroup
type cpuHandler struct {
	prevStats map[string]pair
}

//populateValueList places stats related to cpu handler in in vl. To calculate percentages,
//previous values must be tracked. This function separates previous values according to
//vl.PluginInstance
func (ch *cpuHandler) populateValueList(cpuTime uint64, vl *api.ValueList) {
	if ch.prevStats == nil {
		ch.prevStats = map[string]pair{}
	}

	systemTime := uint64(time.Now().UnixNano())

	vl.Identifier.Type = "pod_cpu"
	vl.DSNames = []string{"percent", "time"}

	var cpuPercent float64
	cpuDelta := float64(cpuTime - ch.prevStats[vl.PluginInstance].cpuTime)
	systemDelta := float64((systemTime - ch.prevStats[vl.PluginInstance].systemTime))

	if cpuDelta > 0.0 && systemDelta > 0.0 && ch.prevStats[vl.PluginInstance] != (pair{}) {
		cpuPercent = (cpuDelta / systemDelta) * 100.0
	}

	ch.prevStats[vl.PluginInstance] = pair{
		systemTime: systemTime,
		cpuTime:    cpuTime,
	}

	vl.Values = []api.Value{api.Gauge(cpuPercent), api.Derive(cpuTime)}
}

//memoryHandler stat gathering and formatting for memory cgroup
type memoryHandler struct{}

//populateValueList places stats related to memory handler in in vl
func (mh *memoryHandler) populateValueList(memStat uint64, vl *api.ValueList) {
	vl.Identifier.Type = "pod_memory"
	vl.Values = []api.Value{api.Gauge(memStat)}
}
