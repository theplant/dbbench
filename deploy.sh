#!/usr/bin/env bash

set -e

server=$1
echo $server

GOOS=linux GOARCH=amd64 go build dynamodb_prepare.go
GOOS=linux GOARCH=amd64 go build dynamodb_query.go

ssh $server -- mkdir -p bench
scp dynamodb_prepare $server:~/bench/
scp dynamodb_query $server:~/bench/
