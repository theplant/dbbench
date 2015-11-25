#!/usr/bin/env bash

export AWS_SECRET_KEY=0
export AWS_SECRET_ACCESS_KEY=0
export AWS_ACCESS_KEY_ID=0
export AWS_ACCESS_KEY=0

go run dynamodb_prepare.go -t 10000 -table hello -g 4 -reset
go run dynamodb_query.go -t 10000 -id 10000 -table hello -g 4

for i in `seq 1 10`; do
	./dynamodb_query -t 10000 -g 1 -v >> query-$i.log 2>&1 &
done
