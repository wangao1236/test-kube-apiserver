#!/bin/bash

IDX=$1

rm ../log/test-curl.log_$IDX
ssh root@127.0.0.1 -p 9996 "rm /var/log/nginx/test.log && \
    systemctl daemon-reload && \
    systemctl restart nginx.service && \
    systemctl status nginx.service"
sleep 3s
go run ../curl/curl.go
cd ../log

cd ../scripts

ssh ao@127.0.0.1 -p 9996 "hostname && \
    cd ~/tmp && \
    cat /var/log/nginx/test.log | grep '10.0.2.2:' > test-curl.log_$IDX && \
    sleep 5s && \
    ls"
scp  -P 9996 ao@127.0.0.1:~/tmp/test-curl.log_$IDX ../log/
go run ../statistics/tcp_statistics.go
ps aux | grep -E "pull-curl" | grep -v "grep" | awk '{print $2}' | xargs kill -9