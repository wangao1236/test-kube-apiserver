package main

import (
	"k8s.io/klog"
	"test-apiserver-client/util"
)

func main()  {
	apiAddrs1, _ := util.ReadFileWithFilter("log/1.67_apiserver.log_6", util.RetryFilter, util.AddrExtractor)
	apiAddrSet1 := util.TransArrayToSet(apiAddrs1)


	apiAddrs2, _ := util.ReadFileWithFilter("log/1.68_apiserver.log_6", util.RetryFilter, util.AddrExtractor)
	apiAddrSet2 := util.TransArrayToSet(apiAddrs2)

	for k1, v1 := range apiAddrSet1 {
		for k2, v2 := range apiAddrSet2 {
			if k1 == k2 {
				klog.Infof("both addr: %+v ,count: %+v %+v", k1, v1, v2)
			}
		}
	}
}
