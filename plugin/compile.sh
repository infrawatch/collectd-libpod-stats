#!/bin/bash

BASEDIR="$HOME/go/src/github.com/pleimer/collectd-libpod-stats"
export COLLECTD_SRC="$HOME/collectd"
export CGO_CPPFLAGS="-I${COLLECTD_SRC}/src/daemon -I${COLLECTD_SRC}/src"
cd "$BASEDIR/plugin"
go build -buildmode=c-shared -o libpodstats.so

rm ../devenv/plugins/libpodstats.so
cp libpodstats.so ../devenv/plugins/

