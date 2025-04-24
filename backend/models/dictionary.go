package models

import (
	"strings"
	utils "wordbuilder/utils"
)

// WordDictionary holds both forward and reverse tries for efficient lookups
type WordDictionary struct {
	WordSet     map[string]bool // For quick word validation
	ForwardTrie *Trie           // For suffix lookups
	ReverseTrie *Trie           // For prefix lookups
	WordList    []string        // Add this field
}

// NewWordDictionary creates a new dictionary with both tries
func NewWordDictionary(wordList []string) *WordDictionary {
	dict := &WordDictionary{
		WordSet:     make(map[string]bool),
		ForwardTrie: NewTrie(),
		ReverseTrie: NewTrie(),
		WordList:    make([]string, 0, len(wordList)),
	}

	for _, word := range wordList {
		word = strings.ToLower(word)
		dict.WordSet[word] = true
		dict.ForwardTrie.Insert(word)
		dict.ReverseTrie.Insert(utils.ReverseString(word))
		dict.WordList = append(dict.WordList, word) // Populate WordList
	}

	return dict
}

// ContainsWord checks if a word exists in the dictionary
func (d *WordDictionary) ContainsWord(word string) bool {
	return d.WordSet[word]
}

// FindWordsWithPrefix returns all words starting with the given prefix
func (d *WordDictionary) FindWordsWithPrefix(prefix string) []string {
	return d.ForwardTrie.KeysWithPrefix(prefix)
}

// FindWordsWithSuffix returns all words ending with the given suffix
func (d *WordDictionary) FindWordsWithSuffix(suffix string) []string {
	reversed := utils.ReverseString(suffix)
	reversedWords := d.ReverseTrie.KeysWithPrefix(reversed)

	// Convert back to normal order
	result := make([]string, len(reversedWords))
	for i, word := range reversedWords {
		result[i] = utils.ReverseString(word)
	}

	return result
}

func (d *WordDictionary) GetForwardTrie() TrieI {
	return d.ForwardTrie
}
func (d *WordDictionary) GetReverseTrie() TrieI {
	return d.ReverseTrie
}
func (d *WordDictionary) GetWordList() []string {
	return d.WordList
}
