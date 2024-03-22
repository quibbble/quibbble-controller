package tsuro

func mapContainsVal(m map[string]string, item string) bool {
	for _, val := range m {
		if val == item {
			return true
		}
	}
	return false
}

func max(m map[string]int) []string {
	currMax := 0
	currKeys := []string{}
	for k, v := range m {
		if v > currMax {
			currMax = v
			currKeys = []string{k}
		} else if v == currMax {
			currKeys = append(currKeys, k)
		}
	}
	return currKeys
}
