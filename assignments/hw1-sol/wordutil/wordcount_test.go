package wordutil

import (
	"testing"
)

func TestWordCountEmpty(t *testing.T) {
	s1 := ""
	m1 := WordCount(s1)
	m1_expected := map[string]int{}
	equals(t, m1_expected, m1)
}

func TestWordCountOneWord(t *testing.T) {
	s1 := "the"
	m1 := WordCount(s1)
	m1_expected := map[string]int{"the": 1}
	equals(t, m1_expected, m1)
}

func TestWordCountTwoWords(t *testing.T) {
	s1 := "the fox"
	m1 := WordCount(s1)
	m1_expected := map[string]int{"the": 1, "fox": 1}
	equals(t, m1_expected, m1)
}

func TestWordCountTwoWordsWithDuplicate(t *testing.T) {
	s1 := "the the"
	m1 := WordCount(s1)
	m1_expected := map[string]int{"the": 2}
	equals(t, m1_expected, m1)
}

func TestWordCountThreeWordsWithDuplicate(t *testing.T) {
	s1 := "the the fox"
	m1 := WordCount(s1)
	m1_expected := map[string]int{"the": 2, "fox": 1}
	equals(t, m1_expected, m1)

	s2 := "the fox the"
	m2 := WordCount(s2)
	m2_expected := map[string]int{"the": 2, "fox": 1}
	equals(t, m2_expected, m2)

	s3 := "the fox fox"
	m3 := WordCount(s3)
	m3_expected := map[string]int{"the": 1, "fox": 2}
	equals(t, m3_expected, m3)
}

func TestWordCountThreeWordsWithDuplicateWithUpperCase(t *testing.T) {
	s1 := "the The fox"
	m1 := WordCount(s1)
	m1_expected := map[string]int{"the": 2, "fox": 1}
	equals(t, m1_expected, m1)

	s2 := "the fox The"
	m2 := WordCount(s2)
	m2_expected := map[string]int{"the": 2, "fox": 1}
	equals(t, m2_expected, m2)

	s3 := "the fox Fox"
	m3 := WordCount(s3)
	m3_expected := map[string]int{"the": 1, "fox": 2}
	equals(t, m3_expected, m3)
}

func TestWordCountMultipleWords(t *testing.T) {
	s1 := "the quick brown fox jumps over the lazy dog"
	m1 := WordCount(s1)
	m1_expected := map[string]int{"the": 2, "fox": 1, "quick": 1, "brown": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}
	equals(t, m1_expected, m1)
}
