#!/bin/bash
go build
./sconn https://localhost:2379 ../certs/ca/ca.crt ../certs/client/client.crt ../certs/client/client.key
