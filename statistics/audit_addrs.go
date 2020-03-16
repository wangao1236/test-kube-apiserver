package main

import (
	"k8s.io/klog"
	"regexp"
	"test-apiserver-client/util"
)

func RequestFilter(log string) bool {
	regex := `requestObject`
	reg, _ := regexp.Compile(regex)
	return reg.FindString(log) != ""
}

func main()  {
	allJobNames, _ := util.ReadFileWithFilter("1.log", util.CreatingFilter, util.JobNameExtractor)

	jobSet := make(map[string]struct{})
	for _, v := range allJobNames {
		jobSet[v] = struct{}{}
	}
	oldNum := len(jobSet)
	klog.Infof("total num: %+v", len(jobSet))

	receiveJobNames1, _ := util.ReadFileWithFilter("log/1.67_audit.log_6", RequestFilter, util.JobNameExtractor)
	receiveJobNameSet1 := util.TransArrayToSet(receiveJobNames1)
	receiveJobNames2, _ := util.ReadFileWithFilter("log/1.68_audit.log_6", RequestFilter, util.JobNameExtractor)
	receiveJobNameSet2 := util.TransArrayToSet(receiveJobNames2)

	for k, v := range receiveJobNameSet1 {
		if _, ok := jobSet[k]; ok {
			delete(jobSet, k)
		}

		if v > 1 {
			klog.Infof("retry job name1: %+v", k)
		}
	}
	for k, v := range receiveJobNameSet2 {
		if _, ok := jobSet[k]; ok {
			delete(jobSet, k)
		}

		if v > 1 {
			klog.Infof("retry job name2: %+v", k)
		}
	}

	bothJobNames := make([]string, 0)
	for k1, _ := range receiveJobNameSet1 {
		for k2, _ := range receiveJobNameSet2 {
			if k1 == k2 {
				bothJobNames = append(bothJobNames, k1)
			}
		}
	}
	klog.Infof("both job names' num: %+v, 1: %+v, 2: %+v", len(bothJobNames), len(receiveJobNameSet1), len(receiveJobNameSet2))

	newNum := len(jobSet)
	klog.Infof("\\'open too many files\\' num: %v", len(jobSet))
	klog.Info(oldNum-newNum)
}
