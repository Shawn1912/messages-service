package utils

import (
	"strings"
	"unicode"
)

// IsPalindrome checks if a given string is a palindrome.
func IsPalindrome(s string) bool {
	var cleaned strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned.WriteRune(unicode.ToLower(r))
		}
	}
	cleanedStr := cleaned.String()
	length := len(cleanedStr)
	for i := 0; i < length/2; i++ {
		if cleanedStr[i] != cleanedStr[length-1-i] {
			return false
		}
	}
	return true
}
