package common

// IncludesString checks if a value exists in a slice of strings
func IncludesString(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}

	return false
}

// FindIndexString returns the index of given value in slice if exists
func FindIndexString(s []string, val string) int {
	for i, v := range s {
		if v == val {
			return i
		}
	}

	return -1
}
