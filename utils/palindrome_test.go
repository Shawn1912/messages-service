package utils

import "testing"

func TestIsPalindrome(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"A man a plan a canal Panama", true},
		{"Hello World", false},
		{"Madam", true},
		{"Not a palindrome", false},
		{"Racecar", true},
		{"", true},
		{" ", true},
		{"12321", true},
		{"12345", false},
		{"Was it a car or a cat I saw", true},
		{"No 'x' in Nixon", true},
		{"No lemon, no melon", true},
		{"Able was I, ere I saw Elba!", true},
		{"Madam, in Eden, I'm Adam.", true},
		{"!@#$", true}, // Only special characters, considered palindrome as cleaned string is empty
	}

	for _, tc := range testCases {
		result := IsPalindrome(tc.input)
		if result != tc.expected {
			t.Errorf("IsPalindrome(%q) = %v; expected %v", tc.input, result, tc.expected)
		}
	}
}
