#!/bin/bash

IDX=$1

rm ~/.kube/config 
scp ao@192.168.1.67:/home/ao/.kube/test-user.config ~/.kube
mv ~/.kube/test-user.config ~/.kube/config
#scp ao@192.168.1.67:/home/ao/.kube/config ~/.kube
ssh root@192.168.1.66 "rm /var/log/nginx/test.log && \
    systemctl daemon-reload && \
    systemctl restart nginx.service && \
    systemctl status nginx.service"

cd ../log || exit
rm 1.67_audit.log_$IDX 1.67_apiserver.log_$IDX 1.68_audit.log_$IDX 1.68_apiserver.log_$IDX nginx.log_$IDX 1.log
cd ../scripts || exit

echo -e "\033[32m ======>>>>>>start send requests \033[0m"
go run ../jobs/batch_create_jobs.go > ../log/1.log 2>&1
echo -e "\033[32m ======>>>>>>finish sending \033[0m"
#nohup go run main.go > 1.log 2>&1 &
#nohup go run main.go > 2.log 2>&1 &
#nohup go run main.go > 3.log 2>&1 &
#nohup go run main.go > 4.log 2>&1 &
#nohup go run main.go > 5.log 2>&1 &
#nohup go run main.go > 6.log 2>&1 &
