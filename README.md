# collectd-libpod-stats
Collectd plugin for gathering resource usage statistics from containers created with the [libpod library](https://github.com/containers/libpod). 

## Current functionality
This plugin is being developed with the new [go-collectd](https://github.com/collectd/go-collectd) library which has not yet released version 1. As a result, the functionality of this project is bound by the ever changing abilities of the library. This means, many features that would be expected from a plugin cannot yet be implemented until the library makes it possible to do so.

As of now, this plugin reports the following for each container on a node:

Resource | Units
---------- | ----------
cpu time | nS
cpu usage | percentage
memory usage | bytes

There is no way to configure these options.

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

