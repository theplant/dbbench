#!/usr/bin/env bash

set -e

server=$1

GOOS=linux GOARCH=amd64 go build dynamodb_prepare.go
GOOS=linux GOARCH=amd64 go build dynamodb_query.go
GOOS=linux GOARCH=amd64 go build postgresql_prepare.go
GOOS=linux GOARCH=amd64 go build postgresql_query.go

ssh $server -- mkdir -p bench
scp dynamodb_prepare $server:~/bench/
scp dynamodb_query $server:~/bench/

# sendagaya338
# theplant
# theplant.cfylmpbyn6qq.us-west-1.rds.amazonaws.com:5432