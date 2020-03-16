#!/bin/bash

#rm ~/.kube/config
#scp -P 9997 ao@127.0.0.1:/home/ao/.kube/test-user.config ~/.kube
#mv ~/.kube/test-user.config ~/.kube/config
rm *.log*
rm ../log/1.67_audit.log_$IDX log/1.67_apiserver.log_$IDX log/1.68_audit.log_$IDX log/1.68_apiserver.log_$IDX log/nginx.log_$IDX
echo -e "\033[32m ======>>>>>>start clear nginx.log \033[0m"
ssh root@127.0.0.1 -p 9996 "rm /var/log/nginx/access.log && \
    ls /var/log/nginx/ && \
    systemctl daemon-reload && \
    systemctl restart nginx.service && \
    systemctl status nginx.service && \
    ls /var/log/nginx/"
echo -e "\033[32m ======>>>>>>finish clearing \033[0m"
echo -e "\033[32m ======>>>>>>start send requests \033[0m"
go run ../jobs/batch_create_jobs.go > 1.log 2>&1
echo -e "\033[32m ======>>>>>>finish sending \033[0m"
#nohup go run main.go > 1.log 2>&1 &
#nohup go run main.go > 2.log 2>&1 &
#nohup go run main.go > 3.log 2>&1 &
#nohup go run main.go > 4.log 2>&1 &
#nohup go run main.go > 5.log 2>&1 &
#nohup go run main.go > 6.log 2>&1 &
