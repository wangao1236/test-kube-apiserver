#!/bin/bash

IDX=$1

cd ../log
rm 1.67_audit.log_$IDX 1.67_apiserver.log_$IDX 1.68_audit.log_$IDX 1.68_apiserver.log_$IDX nginx.log_$IDX
cd ../scripts

ssh ao@127.0.0.1 -p 9996 "hostname && \
    cd ~/tmp && \
    cat /var/log/nginx/access.log | grep '10.0.2.2:' > nginx.log_$IDX && \
    sleep 5s && \
    ls"
scp -P 9996 ao@127.0.0.1:~/tmp/nginx.log_$IDX ../log/

ssh ao@127.0.0.1 -p 9996 "hostname && \
    cd ~/tmp && \
    cat /opt/kubernetes/log/audit.log | grep 'main' > 1.67_audit.log_$IDX && \
    sleep 5s && \
    cat /opt/kubernetes/log/apiserver.log | grep 'main' > 1.67_apiserver.log_$IDX && \
    sleep 5s && \
    ls"
scp -P 9997 ao@127.0.0.1:~/tmp/1.67_audit.log_$IDX ../log/
scp -P 9997 ao@127.0.0.1:~/tmp/1.67_apiserver.log_$IDX ../log/

ssh ao@127.0.0.1:9998 "hostname && \
    cd ~/tmp && \
    cat /opt/kubernetes/log/audit.log | grep 'main' > 1.68_audit.log_$IDX && \
    sleep 5s && \
    cat /opt/kubernetes/log/apiserver.log | grep 'main' > 1.68_apiserver.log_$IDX && \
    sleep 5s && \
    ls"
scp -P 9998 ao@127.0.0.1:~/tmp/1.68_audit.log_$IDX ../log/
scp -P 9998 ao@127.0.0.1:~/tmp/1.68_apiserver.log_$IDX ../log/
