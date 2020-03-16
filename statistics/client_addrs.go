package main

import (
	"k8s.io/klog"
	"test-apiserver-client/util"
)

func main() {
	allJobNames, _ := util.ReadFileWithFilter("1.log", util.CreatingFilter, util.JobNameExtractor)

	jobSet := make(map[string]struct{})
	for _, v := range allJobNames {
		jobSet[v] = struct{}{}
	}
	klog.Info(len(jobSet))

	failedJobNames, _ := util.ReadFileWithFilter("1.log", util.CreateFilter, util.JobNameExtractor)
	for _, v := range failedJobNames {
		if _, ok := jobSet[v]; ok {
			//klog.Info(v)
			delete(jobSet, v)
		}
	}

	//for k, _ := range jobSet {
	//	klog.Infof("failed job's name: %+v", k)
	//}
	klog.Infof("nginx receive requests' num: %+v", len(jobSet))
}
