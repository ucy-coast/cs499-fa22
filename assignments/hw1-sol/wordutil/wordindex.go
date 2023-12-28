package wordutil

import (
	"strings"
)

// Finds first occurrence of each word in a string.
//
// Returns a map that stores each unique word in the string s as the key and
// the index of the first occurence of the word in the input string as the
// corresponding value.
// Matching is case insensitive, e.g. "Orange" and "orange" is considered the
// same word.
func WordIndex(s string) map[string]int {
	var word string
	var word_idx int

	cur := 0
	m := make(map[string]int)

	for len(s) > 0 {
		s_before_trim := s
		s = strings.TrimLeft(s, " ")
		cur = cur + len(s_before_trim) - len(s)
		i := strings.Index(s, " ")
		if i > -1 {
			word = s[:i]
			s = s[i:]
			word_idx = cur
			cur = cur + i
		} else {
			word = s
			word_idx = cur
			s = ""
		}
		word = strings.ToLower(word)
		if _, ok := m[word]; !ok && len(word) > 0 {
			m[word] = word_idx
		}
	}
	return m
}
