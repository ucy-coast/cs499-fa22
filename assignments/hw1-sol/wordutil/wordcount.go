package wordutil

import (
	"strings"
)

// Counts occurrences of each word in a string.
//
// Returns a map that stores each unique word in the string s as the key and
// its count as the corresponding value.
// Matching is case insensitive, e.g. "Orange" and "orange" is considered the
// same word.
func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	m := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(word)
		m[word] += 1
	}
	return m
}
