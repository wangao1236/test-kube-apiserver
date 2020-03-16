package main

import (
	"k8s.io/klog"
	"test-apiserver-client/util"
)

func main()  {
	allJobNames, _ := util.ReadFileWithFilter("1.log", util.CreatingFilter, util.JobNameExtractor)

	jobSet := make(map[string]struct{})
	for _, v := range allJobNames {
		jobSet[v] = struct{}{}
	}
	klog.Info(len(jobSet))

	nginxJobNames, _ := util.ReadFileWithFilter("log/nginx.log_6", NoFilter, util.JobNameExtractor)

	for _, v := range nginxJobNames {
		if _, ok := jobSet[v]; ok {
			delete(jobSet, v)
		}
	}

	//for k, _ := range jobSet {
	//	klog.Infof("failed job's name: %+v", k)
	//}
	klog.Infof("\\'open too many files\\' num: %v", len(jobSet))
}
