package models

import (
	"testing"
)

func TestNewWordDictionaryAndContainsWord(t *testing.T) {
	words := []string{"apple", "banana", "grape"}
	dict := NewWordDictionary(words)

	for _, w := range words {
		if !dict.ContainsWord(w) {
			t.Errorf("Expected %s to be found", w)
		}
	}
	if dict.ContainsWord("orange") {
		t.Errorf("Did not expect 'orange' to be found")
	}
}

func TestFindWordsWithPrefix(t *testing.T) {
	words := []string{"apple", "applet", "banana"}
	dict := NewWordDictionary(words)
	results := dict.FindWordsWithPrefix("app")
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestFindWordsWithSuffix(t *testing.T) {
	words := []string{"testing", "sing", "bring"}
	dict := NewWordDictionary(words)
	results := dict.FindWordsWithSuffix("ing")
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
}
