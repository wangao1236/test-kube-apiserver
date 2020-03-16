#!/bin/bash

IDX=$1

cd ../log || exit
rm 1.67_audit.log_$IDX 1.67_apiserver.log_$IDX 1.68_audit.log_$IDX 1.68_apiserver.log_$IDX nginx.log_$IDX
cd ../scripts || exit

pwd

ssh ao@192.168.1.66 "hostname && \
    cd ~/tmp && \
    cat /var/log/nginx/access.log | grep '192.168.1.6:' > nginx.log_$IDX && \
    sleep 5s && \
    ls"
scp ao@192.168.1.66:~/tmp/nginx.log_$IDX ../log/

ssh ao@192.168.1.67 "hostname && \
    cd ~/tmp && \
    cat /opt/kubernetes/log/audit.log | grep 'batch_create_jobs' > 1.67_audit.log_$IDX && \
    sleep 5s && \
    cat /opt/kubernetes/log/apiserver.log | grep 'batch_create_jobs' > 1.67_apiserver.log_$IDX && \
    sleep 5s && \
    ls"
scp ao@192.168.1.67:~/tmp/1.67_audit.log_$IDX ../log/
scp ao@192.168.1.67:~/tmp/1.67_apiserver.log_$IDX ../log/

ssh ao@192.168.1.68 "hostname && \
    cd ~/tmp && \
    cat /opt/kubernetes/log/audit.log | grep 'batch_create_jobs' > 1.68_audit.log_$IDX && \
    sleep 5s && \
    cat /opt/kubernetes/log/apiserver.log | grep 'batch_create_jobs' > 1.68_apiserver.log_$IDX && \
    sleep 5s && \
    ls"
scp ao@192.168.1.68:~/tmp/1.68_audit.log_$IDX ../log/
scp ao@192.168.1.68:~/tmp/1.68_apiserver.log_$IDX ../log/
