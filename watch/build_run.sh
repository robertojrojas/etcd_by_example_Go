#!/bin/bash
rm -rf default.etcd
rm -rf watch
go build
key=${1:-golang}
./watch ${key}
