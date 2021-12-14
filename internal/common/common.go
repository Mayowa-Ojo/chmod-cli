package common

import "golang.design/x/clipboard"

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

// CopyToClipboard copies string to clipboard
func CopyToClipboard(s string) {
	clipboard.Write(clipboard.FmtText, []byte(s))
}
