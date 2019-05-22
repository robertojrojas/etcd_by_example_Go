#!/bin/bash
rm -rf default.etcd
rm -rf embed
go build
./embed http://localhost:2379
