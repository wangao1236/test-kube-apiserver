package util

func TransArrayToSet(old []string) map[string]int  {
	//klog.Info(len(old))
	result := make(map[string]int)
	for _, v := range old {
		if _, ok := result[v] ; !ok {
			result[v] = 1
		} else {
			result[v]++
		}
	}
	//klog.Info(len(result))
	return result
}
