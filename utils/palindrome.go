package utils

import "strings"

func IsPalindrome(s string) bool {
	sanitized := strings.ToLower(strings.ReplaceAll(s, " ", ""))
	length := len(sanitized)
	for i := 0; i < length/2; i++ {
		if sanitized[i] != sanitized[length-1-i] {
			return false
		}
	}
	return true
}
