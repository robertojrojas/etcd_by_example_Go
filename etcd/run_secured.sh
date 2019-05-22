#!/bin/bash

etcd  --client-cert-auth --trusted-ca-file=../certs/ca/ca.crt --listen-client-urls=https://10.0.1.198:2379,https://localhost:2379 --advertise-client-urls=https://10.0.1.198:2379 --cert-file=../certs/server/server.crt --key-file=../certs/server/server.key
