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
