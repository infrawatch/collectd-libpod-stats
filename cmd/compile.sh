#!/bin/bash

BASEDIR="$HOME/go/src/github.com/pleimer/collectd-podman-stats"
export COLLECTD_SRC="$HOME/collectd"
export CGO_CPPFLAGS="-I${COLLECTD_SRC}/src/daemon -I${COLLECTD_SRC}/src"
cd "$BASEDIR/cmd"
go build -buildmode=c-shared -o podmanstats.so

rm ../devenv/plugins/podmanstats.so
cp podmanstats.so ../devenv/plugins/

