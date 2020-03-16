package util

import "regexp"

func JobNameExtractor(log string) string {
	// test-job-284-1583933142133706000-acde48001122
	regex := `test-job-\d+-\d+-\S{12}`
	reg, _ := regexp.Compile(regex)
	return reg.FindString(log)
}

func AddrExtractor(log string) string {
	// test-job-284-1583933142133706000-acde48001122
	regex := `(\d{1,3}\.){3}\d{1,3}:\d{1,5}`
	reg, _ := regexp.Compile(regex)
	return reg.FindString(log)
}
