package utils

import (
	"strings"
	"unicode"
)

// IsPalindrome checks if a given string is a palindrome.
func IsPalindrome(s string) bool {
	// Mapping function to clean the string
	f := func(r rune) rune {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return -1 // Exclude the character
		}
		return unicode.ToLower(r)
	}
	// Clean the string
	cleaned := strings.Map(f, s)

	// Convert to a slice of runes to handle Unicode characters
	runes := []rune(cleaned)
	i, j := 0, len(runes)-1

	// Compare characters from both ends
	for i < j {
		if runes[i] != runes[j] {
			return false
		}
		i++
		j--
	}
	return true
}
