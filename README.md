# collectd-libpod-stats
Collectd plugin for gathering resource usage statistics from containers created with the [libpod library](https://github.com/containers/libpod). 

## Functionality
This plugin is being developed with the new [go-collectd](https://github.com/collectd/go-collectd) library which has not yet released version 1. As a result, the functionality of this project is bound by the ever changing abilities of the library. Many features that would be expected from a plugin cannot yet be implemented until the library makes it possible to do so.

Right now, this plugin reports the following for each container on a node:

Resource | Units
---------- | ----------
cpu time | nS
cpu usage | percentage
memory usage | bytes

There is no way to configure these options.

## Notes on memory usage calculations
Calculating memory usage on linux systems is hairy; not all projects do it the same way. This plugin replicates the memory calculation that the podman tool provides, using the value of `memory.usage_in_bytes` for cgroups v1 or `memory.current` for cgroups v2. This value includes both RSS and CACHE memory and is different to the value returned by the `free` tool, which calculates used memory by reading `/proc/meminfo`.

Kubernetes calulates the memory usage as `memory.usage_in_bytes - total_inactive_file` [(source)](https://github.com/kubernetes/kubernetes/blob/dde6e8e7465468c32642659cb708a5cc922add64/test/e2e/node/node_problem_detector.go#L242). K8s uses the value to make [eviction decisions](https://kubernetes.io/docs/tasks/administer-cluster/out-of-resource/#eviction-signals). Because k8s provides its own meanss of monitoring pod memory and calculates usage differently than this plugin and the `free` tool, this plugin should not be used to monitor pods in a k8s cluster.

Ultimately, the stated values are used here because the linux cgroup oom_killer uses this value to destroy processes that exceed the hard limit.

## Build from source

Collectd library files must be made available to the golang compiler to build this plugin. Clone the [collectd project](https://github.com/collectd/collectd).
Then, run the following:

```bash
export COLLECTD_SRC="/path/to/collectd/source"
export CGO_CPPFLAGS="-I${COLLECTD_SRC}/src/daemon -I${COLLECTD_SRC}/src"


git clone https://github.com/infrawatch/collectd-libpod-stats.git
cd collectd-libpod-stats/plugin
go build -buildmode=c-shared -o libpodstats.so
```

## collectd.conf

```xml
TypesDB "/path/to/types.db.libpodstats"

LoadPlugin "libpodstats"
<Plugin "libpodstats">
</Plugin>
```

