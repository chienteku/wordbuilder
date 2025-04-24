package models

import (
	"sort"
	"testing"
)

func TestTrie_InsertAndContains(t *testing.T) {
	trie := NewTrie()
	words := []string{"apple", "app", "banana", "band", "bandana", "cab"}

	// Insert words
	for _, word := range words {
		trie.Insert(word)
	}

	// Should contain all inserted words
	for _, word := range words {
		if !trie.Contains(word) {
			t.Errorf("Trie should contain %q", word)
		}
	}

	// Should not contain words not inserted
	notWords := []string{"appl", "ban", "banda", "ca", "dog"}
	for _, word := range notWords {
		if trie.Contains(word) {
			t.Errorf("Trie should NOT contain %q", word)
		}
	}
}

func TestTrie_KeysWithPrefix(t *testing.T) {
	trie := NewTrie()
	words := []string{"apple", "app", "banana", "band", "bandana", "cab"}
	for _, word := range words {
		trie.Insert(word)
	}

	tests := []struct {
		prefix   string
		expected []string
	}{
		{"app", []string{"app", "apple"}},
		{"ban", []string{"banana", "band", "bandana"}},
		{"band", []string{"band", "bandana"}},
		{"c", []string{"cab"}},
		{"z", []string{}},
		{"", []string{"app", "apple", "banana", "band", "bandana", "cab"}},
	}

	for _, tt := range tests {
		got := trie.KeysWithPrefix(tt.prefix)
		sort.Strings(got)
		sort.Strings(tt.expected)
		if !equalStringSlices(got, tt.expected) {
			t.Errorf("KeysWithPrefix(%q) = %v; want %v", tt.prefix, got, tt.expected)
		}
	}
}

func TestTrie_GetNextLetters(t *testing.T) {
	trie := NewTrie()
	words := []string{"apple", "app", "banana", "band", "bandana", "cab"}
	for _, word := range words {
		trie.Insert(word)
	}

	tests := []struct {
		prefix   string
		expected []string
	}{
		{"", []string{"a", "b", "c"}},
		{"app", []string{"l"}},
		{"ban", []string{"a", "d"}},
		{"band", []string{"a"}},
		{"bandana", []string{}},
		{"z", []string{}},
	}

	for _, tt := range tests {
		got := trie.GetNextLetters(tt.prefix)
		sort.Strings(got)
		sort.Strings(tt.expected)
		if !equalStringSlices(got, tt.expected) {
			t.Errorf("GetNextLetters(%q) = %v; want %v", tt.prefix, got, tt.expected)
		}
	}
}
