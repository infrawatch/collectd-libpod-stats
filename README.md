# collectd-libpod-stats
Collectd plugin for gathering resource usage statistics from containers created with the [libpod library](https://github.com/containers/libpod). 

## Functionality
This plugin is being developed with the new [go-collectd](https://github.com/collectd/go-collectd) library which has not yet released version 1. As a result, the functionality of this project is bound by the ever changing abilities of the library. Many features that would be expected from a plugin cannot yet be implemented until the library makes it possible to do so.

As of now, this plugin reports the following for each container on a node:

Resource | Units
---------- | ----------
cpu time | nS
cpu usage | percentage
memory usage | bytes

There is no way to configure these options.

## Notes on memory usage calculations
Calculating memory usage on linux systems is hairy and differs based on projects. This plugin replicates the memory calculation that the podman tool uses, that is, the value of `memory.usage_in_bytes` in the case of cgroups v1 or `memory.current` for cgroups v2. This value includes both RSS and CACHE memory values and thus is different to what the `free` tool reports as used memory by reading `/proc/meminfo`.

Kubernetes calulates the memory usage as `memory.usage_in_bytes - total_inactive_file` [(source)](https://github.com/kubernetes/kubernetes/blob/dde6e8e7465468c32642659cb708a5cc922add64/test/e2e/node/node_problem_detector.go#L242). Kubernetes uses this value to make [eviction decisions](https://kubernetes.io/docs/tasks/administer-cluster/out-of-resource/#eviction-signals) and thus this value is important. However, k8s provides its own means of monitoring pod memory and should be used instead of this plugin.

Ultimately, this plugin uses the stated values because the linux cgroup oom_killer uses this value to destroy processes that exceed the hard limit.

## Build from source

Collectd library files must be made available to the golang compiler to build this plugin. Clone the [collectd project](https://github.com/collectd/collectd).
Then, run the following:

```bash
export COLLECTD_SRC="/path/to/collectd/source"
export CGO_CPPFLAGS="-I${COLLECTD_SRC}/src/daemon -I${COLLECTD_SRC}/src"


git clone https://github.com/pleimer/collectd-libpod-stats.git
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

