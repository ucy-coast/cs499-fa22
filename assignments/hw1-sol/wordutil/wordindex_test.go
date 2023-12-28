package wordutil

import (
	"testing"
)

func TestWordIndexEmpty(t *testing.T) {
	s1 := ""
	m1 := WordIndex(s1)
	m1_expected := map[string]int{}
	equals(t, m1_expected, m1)
}

func TestWordIndexOneWord(t *testing.T) {
	s1 := "the"
	m1 := WordIndex(s1)
	m1_expected := map[string]int{"the": 0}
	equals(t, m1_expected, m1)
}

func TestWordIndexTwoWords(t *testing.T) {
	s1 := "the fox"
	m1 := WordIndex(s1)
	m1_expected := map[string]int{"the": 0, "fox": 4}
	equals(t, m1_expected, m1)
}

func TestWordIndexThreeWordsWithDuplicate(t *testing.T) {
	s1 := "the the fox"
	m1 := WordIndex(s1)
	m1_expected := map[string]int{"the": 0, "fox": 8}
	equals(t, m1_expected, m1)

	s2 := "the fox the"
	m2 := WordIndex(s2)
	m2_expected := map[string]int{"the": 0, "fox": 4}
	equals(t, m2_expected, m2)

	s3 := "the fox fox"
	m3 := WordIndex(s3)
	m3_expected := map[string]int{"the": 0, "fox": 4}
	equals(t, m3_expected, m3)
}

func TestWordIndexThreeWordsWithDuplicateWithUpperCase(t *testing.T) {
	s1 := "the The fox"
	m1 := WordIndex(s1)
	m1_expected := map[string]int{"the": 0, "fox": 8}
	equals(t, m1_expected, m1)

	s2 := "the fox The"
	m2 := WordIndex(s2)
	m2_expected := map[string]int{"the": 0, "fox": 4}
	equals(t, m2_expected, m2)

	s3 := "the fox Fox"
	m3 := WordIndex(s3)
	m3_expected := map[string]int{"the": 0, "fox": 4}
	equals(t, m3_expected, m3)
}

func TestWordIndexMultipleWords(t *testing.T) {
	s1 := "the quick brown fox jumps over the lazy dog"
	m1 := WordIndex(s1)
	m1_expected := map[string]int{"the": 0, "quick": 4, "brown": 10, "fox": 16, "jumps": 20, "over": 26, "lazy": 35, "dog": 40}
	equals(t, m1_expected, m1)
}
