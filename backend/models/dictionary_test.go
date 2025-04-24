package models

import (
	"sort"
	"testing"
)

// Sample word list for testing
var sampleWords = []string{"apple", "banana", "band", "bandana", "cab", "can"}

func newTestDictionary() *WordDictionary {
	return NewWordDictionary(sampleWords)
}

func TestNewWordDictionary(t *testing.T) {
	dict := newTestDictionary()
	if len(dict.WordSet) != len(sampleWords) {
		t.Errorf("Expected WordSet length %d, got %d", len(sampleWords), len(dict.WordSet))
	}
	for _, word := range sampleWords {
		if !dict.WordSet[word] {
			t.Errorf("Word %q not found in WordSet", word)
		}
	}
}

func TestContainsWord(t *testing.T) {
	dict := newTestDictionary()
	for _, word := range sampleWords {
		if !dict.ContainsWord(word) {
			t.Errorf("Expected ContainsWord(%q) to be true", word)
		}
	}
	// Test for a word not in the dictionary
	if dict.ContainsWord("dog") {
		t.Errorf("Expected ContainsWord(\"dog\") to be false")
	}
}

func TestFindWordsWithPrefix(t *testing.T) {
	dict := newTestDictionary()
	tests := []struct {
		prefix   string
		expected []string
	}{
		{"ban", []string{"banana", "band", "bandana"}},
		{"ba", []string{"banana", "band", "bandana"}},
		{"b", []string{"banana", "band", "bandana"}},
		{"c", []string{"cab", "can"}},
		{"a", []string{"apple"}},
		{"z", []string{}},
	}
	for _, tt := range tests {
		results := dict.FindWordsWithPrefix(tt.prefix)
		sort.Strings(results)
		sort.Strings(tt.expected)
		if !equalStringSlices(results, tt.expected) {
			t.Errorf("FindWordsWithPrefix(%q) = %v; want %v", tt.prefix, results, tt.expected)
		}
	}
}

func TestFindWordsWithSuffix(t *testing.T) {
	dict := newTestDictionary()
	tests := []struct {
		suffix   string
		expected []string
	}{
		{"ana", []string{"banana", "bandana"}},
		{"le", []string{"apple"}},
		{"d", []string{"band"}},
		{"n", []string{"can"}},
		{"z", []string{}},
	}
	for _, tt := range tests {
		results := dict.FindWordsWithSuffix(tt.suffix)
		sort.Strings(results)
		sort.Strings(tt.expected)
		if !equalStringSlices(results, tt.expected) {
			t.Errorf("FindWordsWithSuffix(%q) = %v; want %v", tt.suffix, results, tt.expected)
		}
	}
}

func TestGetForwardTrie_ReverseTrie_WordList(t *testing.T) {
	dict := newTestDictionary()
	if dict.GetForwardTrie() == nil {
		t.Error("GetForwardTrie returned nil")
	}
	if dict.GetReverseTrie() == nil {
		t.Error("GetReverseTrie returned nil")
	}
	wordList := dict.GetWordList()
	sort.Strings(wordList)
	sort.Strings(sampleWords)
	if !equalStringSlices(wordList, sampleWords) {
		t.Errorf("GetWordList() = %v; want %v", wordList, sampleWords)
	}
}
