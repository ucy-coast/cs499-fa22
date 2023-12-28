package wordutil

import (
	"testing"
)

func TestWordIndexAllEmpty(t *testing.T) {
	s1 := ""
	m1 := WordIndexAll(s1)
	m1_expected := map[string][]int{}
	equals(t, m1_expected, m1)
}

func TestWordIndexAllOneWord(t *testing.T) {
	s1 := "the"
	m1 := WordIndexAll(s1)
	m1_expected := map[string][]int{"the": {0}}
	equals(t, m1_expected, m1)
}

func TestWordIndexAllTwoWords(t *testing.T) {
	s1 := "the fox"
	m1 := WordIndexAll(s1)
	m1_expected := map[string][]int{"the": {0}, "fox": {4}}
	equals(t, m1_expected, m1)

	s2 := "t f"
	m2 := WordIndexAll(s2)
	m2_expected := map[string][]int{"t": {0}, "f": {2}}
	equals(t, m2_expected, m2)
}

func TestWordIndexAllThreeWordsWithDuplicate(t *testing.T) {
	s1 := "the the fox"
	m1 := WordIndexAll(s1)
	m1_expected := map[string][]int{"the": {0, 4}, "fox": {8}}
	equals(t, m1_expected, m1)

	s2 := "the fox the"
	m2 := WordIndexAll(s2)
	m2_expected := map[string][]int{"the": {0, 8}, "fox": {4}}
	equals(t, m2_expected, m2)

	s3 := "the fox fox"
	m3 := WordIndexAll(s3)
	m3_expected := map[string][]int{"the": {0}, "fox": {4, 8}}
	equals(t, m3_expected, m3)
}

func TestWordIndexAllThreeWordsWithDuplicateWithUpperCase(t *testing.T) {
	s1 := "the The fox"
	m1 := WordIndexAll(s1)
	m1_expected := map[string][]int{"the": {0, 4}, "fox": {8}}
	equals(t, m1_expected, m1)

	s2 := "the fox The"
	m2 := WordIndexAll(s2)
	m2_expected := map[string][]int{"the": {0, 8}, "fox": {4}}
	equals(t, m2_expected, m2)

	s3 := "the fox Fox"
	m3 := WordIndexAll(s3)
	m3_expected := map[string][]int{"the": {0}, "fox": {4, 8}}
	equals(t, m3_expected, m3)
}

func TestWordIndexAllMultipleWords(t *testing.T) {
	s1 := "the quick brown fox jumps over the lazy dog"
	m1 := WordIndexAll(s1)
	m1_expected := map[string][]int{"the": {0, 31}, "quick": {4}, "brown": {10}, "fox": {16}, "jumps": {20}, "over": {26}, "lazy": {35}, "dog": {40}}
	equals(t, m1_expected, m1)
}
