package util

import (
	"regexp"
	"strings"
)

func CreatingFilter(log string) bool {
	return strings.Contains(log, "Creating job")
}

func CreateFilter(log string) bool {
	return strings.Contains(log, "Create job")
}

func CreatedFilter(log string) bool {
	regex := `Created job`
	reg, _ := regexp.Compile(regex)
	return reg.FindString(log) != ""
}


func NoFilter(log string) bool {
	return true
}

func RetryFilter(log string) bool {
	regex := `429 \[main\/v0.0.0 \(darwin\/amd64\)`
	reg, _ := regexp.Compile(regex)
	//if reg.MatchString(log) {
	//	klog.Info(log)
	//}
	return reg.MatchString(log)
}

func Status429Filter(log string) bool {
	regex := `\"ups_status\":\"429, 201\"`
	reg, _ := regexp.Compile(regex)
	//if reg.MatchString(log) {
	//	klog.Info(log)
	//}
	return reg.MatchString(log)
}
