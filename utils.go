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
