package main

import (
	"k8s.io/klog"
	"test-apiserver-client/util"
)

func main()  {
	cnt, _ := util.ReadWriteFileWithFilter("log/nginx.log_6", "log/429.log_6", util.Status429Filter)

	klog.Infof("nginx 429 cnt: %+v", cnt)

	apiAddrs1, _ := util.ReadFileWithFilter("log/1.67_apiserver.log_6", util.RetryFilter, util.AddrExtractor)
	apiAddrSet1 := util.TransArrayToSet(apiAddrs1)


	apiAddrs2, _ := util.ReadFileWithFilter("log/1.68_apiserver.log_6", util.RetryFilter, util.AddrExtractor)
	apiAddrSet2 := util.TransArrayToSet(apiAddrs2)

	klog.Infof("apiserver 429 cnt: %+v %+v", len(apiAddrSet1), len(apiAddrSet2))
}
