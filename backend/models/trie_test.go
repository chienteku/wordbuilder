package models

import "testing"

func TestTrieInsertAndContains(t *testing.T) {
	trie := NewTrie()
	trie.Insert("cat")
	trie.Insert("car")

	if !trie.Contains("cat") || !trie.Contains("car") {
		t.Error("Expected words to be found in trie")
	}
	if trie.Contains("dog") {
		t.Error("Did not expect 'dog' to be found")
	}
}

func TestTrieKeysWithPrefix(t *testing.T) {
	trie := NewTrie()
	trie.Insert("cat")
	trie.Insert("car")
	trie.Insert("dog")

	results := trie.KeysWithPrefix("ca")
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestTrieGetNextLetters(t *testing.T) {
	trie := NewTrie()
	trie.Insert("cat")
	trie.Insert("car")
	letters := trie.GetNextLetters("ca")
	if len(letters) != 2 {
		t.Errorf("Expected 2 next letters, got %d", len(letters))
	}
}
