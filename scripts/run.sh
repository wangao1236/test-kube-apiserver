#!/bin/bash

IDX=$1

./test.sh $IDX
#sleep 60s
ps aux | grep -E "exe/batch_create_jobs|batch_create_jobs.go" | grep -v "grep" | awk '{print $2}' | xargs kill -9
sleep 3s
./pull.sh $IDX

go run ../statistics/tcp_statistics.go
