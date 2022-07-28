package go_tsuro

func contains(list []string, item string) bool {
	for _, val := range list {
		if val == item {
			return true
		}
	}
	return false
}

func mapContainsVal(m map[string]string, item string) bool {
	for _, val := range m {
		if val == item {
			return true
		}
	}
	return false
}

func indexOf(items []string, item string) int {
	for index, it := range items {
		if it == item {
			return index
		}
	}
	return -1
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
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

func duplicates(list []string) bool {
	for idx, v1 := range list {
		for _, v2 := range list[idx+1:] {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}
