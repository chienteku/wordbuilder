package main

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// TrieNode represents a node in the Trie
type TrieNode struct {
	Children map[rune]*TrieNode
	IsWord   bool
}

// Trie represents a dictionary structure for efficient prefix searches
type Trie struct {
	Root *TrieNode
}

// NewTrie creates a new Trie
func NewTrie() *Trie {
	return &Trie{Root: &TrieNode{Children: make(map[rune]*TrieNode)}}
}

// Insert adds a word to the Trie
func (t *Trie) Insert(word string) {
	node := t.Root
	for _, ch := range word {
		if _, exists := node.Children[ch]; !exists {
			node.Children[ch] = &TrieNode{Children: make(map[rune]*TrieNode)}
		}
		node = node.Children[ch]
	}
	node.IsWord = true
}

// Contains checks if a word exists in the Trie
func (t *Trie) Contains(word string) bool {
	node := t.Root
	for _, ch := range word {
		if _, exists := node.Children[ch]; !exists {
			return false
		}
		node = node.Children[ch]
	}
	return node.IsWord
}

// KeysWithPrefix returns all words that start with the given prefix
func (t *Trie) KeysWithPrefix(prefix string) []string {
	var results []string
	node := t.Root
	for _, ch := range prefix {
		if _, exists := node.Children[ch]; !exists {
			return results
		}
		node = node.Children[ch]
	}
	t.collectKeys(node, prefix, &results)
	return results
}

// collectKeys recursively collects all words from the current node
func (t *Trie) collectKeys(node *TrieNode, prefix string, results *[]string) {
	if node.IsWord {
		*results = append(*results, prefix)
	}
	for ch, child := range node.Children {
		t.collectKeys(child, prefix+string(ch), results)
	}
}

// GetNextLetters returns the possible next letters after a given prefix
func (t *Trie) GetNextLetters(prefix string) []string {
	node := t.Root
	for _, ch := range prefix {
		if child, exists := node.Children[ch]; exists {
			node = child
		} else {
			return []string{}
		}
	}
	letters := make([]string, 0, len(node.Children))
	for ch := range node.Children {
		letters = append(letters, string(ch))
	}
	return letters
}

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
		dict.ReverseTrie.Insert(reverseString(word))
		dict.WordList = append(dict.WordList, word) // Populate WordList
	}

	return dict
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
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
	reversed := reverseString(suffix)
	reversedWords := d.ReverseTrie.KeysWithPrefix(reversed)

	// Convert back to normal order
	result := make([]string, len(reversedWords))
	for i, word := range reversedWords {
		result[i] = reverseString(word)
	}

	return result
}

// EnhancedWordBuilder represents the word building game state
type EnhancedWordBuilder struct {
	Answer      string
	PrefixSet   map[string]bool
	SuffixSet   map[string]bool
	Step        int
	Dictionary  *WordDictionary
	IsValidWord bool

	// Enhanced features
	ValidCompletions []string // Complete valid words that can be formed
	Suggestion       string   // Suggested next move
}

// NewEnhancedWordBuilder creates a new word builder instance
func NewEnhancedWordBuilder(dictionary *WordDictionary) *EnhancedWordBuilder {
	wb := &EnhancedWordBuilder{
		Answer:           "",
		PrefixSet:        make(map[string]bool),
		SuffixSet:        make(map[string]bool),
		Step:             0,
		Dictionary:       dictionary,
		IsValidWord:      false,
		ValidCompletions: []string{},
		Suggestion:       "",
	}

	wb.UpdateSets()
	return wb
}

// CheckValidWord verifies if the current answer is a valid word
func (wb *EnhancedWordBuilder) CheckValidWord() bool {
	isValid := len(wb.Answer) > 0 && wb.Dictionary.ContainsWord(wb.Answer)
	wb.IsValidWord = isValid
	return isValid
}

// Reset clears the current state
func (wb *EnhancedWordBuilder) Reset() {
	wb.Answer = ""
	wb.Step = 0
	wb.IsValidWord = false
	wb.ValidCompletions = []string{}
	wb.Suggestion = ""
	wb.UpdateSets()
}

// AddLetter adds a letter to either the prefix or suffix
func (wb *EnhancedWordBuilder) AddLetter(letter, position string) (bool, string) {
	letter = strings.ToLower(letter)

	if position == "prefix" {
		if !wb.PrefixSet[letter] {
			return false, fmt.Sprintf("Invalid letter '%s' for prefix position.", letter)
		}
		wb.Answer = letter + wb.Answer
	} else if position == "suffix" {
		if !wb.SuffixSet[letter] {
			return false, fmt.Sprintf("Invalid letter '%s' for suffix position.", letter)
		}
		wb.Answer = wb.Answer + letter
	} else {
		return false, "Invalid position. Use 'prefix' or 'suffix'."
	}

	wb.CheckValidWord()
	wb.UpdateSets()
	wb.Step++

	message := fmt.Sprintf("Step %d: Added '%s' as %s -> Answer: %s", wb.Step, letter, position, wb.Answer)
	if wb.IsValidWord {
		message += fmt.Sprintf("\n*** '%s' is a valid word! ***", wb.Answer)
	}

	// Add information about possible completions
	if len(wb.ValidCompletions) > 0 && !wb.IsValidWord {
		message += fmt.Sprintf("\nPossible completions: %s", strings.Join(wb.ValidCompletions[:min(3, len(wb.ValidCompletions))], ", "))
	}

	return true, message
}

// RemoveLetter removes a letter at the specified index
func (wb *EnhancedWordBuilder) RemoveLetter(index int) (bool, string) {
	if index < 0 || index >= len(wb.Answer) {
		return false, fmt.Sprintf("Invalid index %d for answer '%s'.", index, wb.Answer)
	}

	letter := string(wb.Answer[index])
	newAnswer := wb.Answer[:index] + wb.Answer[index+1:]
	wb.Answer = newAnswer

	wb.CheckValidWord()
	wb.UpdateSets()
	wb.Step++

	message := fmt.Sprintf("Step %d: Removed '%s' at index %d -> Answer: %s", wb.Step, letter, index, wb.Answer)
	if wb.IsValidWord {
		message += fmt.Sprintf("\n*** '%s' is a valid word! ***", wb.Answer)
	}

	return true, message
}

// UpdateSets updates the available prefix and suffix letter sets
func (wb *EnhancedWordBuilder) UpdateSets() {
	wb.PrefixSet = make(map[string]bool)
	wb.SuffixSet = make(map[string]bool)
	wb.ValidCompletions = []string{}

	// If no letters yet, provide all letters that can start or end words
	if len(wb.Answer) == 0 {
		for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
			letterStr := string(letter)
			// Check if this letter can start any words
			if len(wb.Dictionary.FindWordsWithPrefix(letterStr)) > 0 {
				wb.SuffixSet[letterStr] = true
			}
			// Check if this letter can end any words
			if len(wb.Dictionary.FindWordsWithSuffix(letterStr)) > 0 {
				wb.PrefixSet[letterStr] = true
			}
		}
		return
	}

	foundValidContinuation := false

	// 1. Suffix letters from ForwardTrie
	suffixLetters := wb.Dictionary.ForwardTrie.GetNextLetters(wb.Answer)
	for _, letter := range suffixLetters {
		wb.SuffixSet[letter] = true
		foundValidContinuation = true
	}
	// Optionally collect some completions
	if len(suffixLetters) > 0 {
		words := wb.Dictionary.FindWordsWithPrefix(wb.Answer)
		for i, word := range words {
			if i >= 5 { // Limit to reduce overhead
				break
			}
			if len(word) > len(wb.Answer) {
				wb.ValidCompletions = append(wb.ValidCompletions, word)
			}
		}
	}

	// 2. Prefix letters from ReverseTrie
	reversedAnswer := reverseString(wb.Answer)
	prefixLetters := wb.Dictionary.ReverseTrie.GetNextLetters(reversedAnswer)
	for _, letter := range prefixLetters {
		wb.PrefixSet[letter] = true
		foundValidContinuation = true
	}
	// Optionally collect some completions
	if len(prefixLetters) > 0 {
		revWords := wb.Dictionary.ReverseTrie.KeysWithPrefix(reversedAnswer)
		for i, revWord := range revWords {
			if i >= 5 { // Limit to reduce overhead
				break
			}
			if len(revWord) > len(reversedAnswer) {
				wb.ValidCompletions = append(wb.ValidCompletions, reverseString(revWord))
			}
		}
	}

	// 3. Embedded check with parallelism
	wordList := wb.Dictionary.WordList
	numParts := runtime.NumCPU()
	partSize := (len(wordList) + numParts - 1) / numParts
	var wg sync.WaitGroup
	type Result struct {
		prefixSet map[string]bool
		suffixSet map[string]bool
	}
	results := make(chan Result, numParts)

	for i := 0; i < numParts; i++ {
		wg.Add(1)
		start := i * partSize
		end := start + partSize
		if end > len(wordList) {
			end = len(wordList)
		}
		go func(words []string) {
			defer wg.Done()
			localPrefix := make(map[string]bool)
			localSuffix := make(map[string]bool)
			for _, word := range words {
				idx := strings.Index(word, wb.Answer)
				if idx >= 0 {
					foundValidContinuation = true
					if idx > 0 {
						localPrefix[string(word[idx-1])] = true
					}
					embedEndIdx := idx + len(wb.Answer)
					if embedEndIdx < len(word) {
						localSuffix[string(word[embedEndIdx])] = true
					}
				}
			}
			results <- Result{localPrefix, localSuffix}
		}(wordList[start:end])
	}

	// Collect results concurrently
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		for letter := range res.prefixSet {
			wb.PrefixSet[letter] = true
		}
		for letter := range res.suffixSet {
			wb.SuffixSet[letter] = true
		}
	}

	// Generate suggestion (simplified)
	if len(wb.ValidCompletions) > 0 && !wb.IsValidWord {
		sort.Slice(wb.ValidCompletions, func(i, j int) bool {
			return len(wb.ValidCompletions[i]) < len(wb.ValidCompletions[j])
		})
		shortest := wb.ValidCompletions[0]
		if strings.HasPrefix(shortest, wb.Answer) && len(shortest) > len(wb.Answer) {
			wb.Suggestion = "Try adding '" + string(shortest[len(wb.Answer)]) + "' as suffix"
		} else if idx := strings.Index(shortest, wb.Answer); idx > 0 {
			wb.Suggestion = "Try adding '" + string(shortest[idx-1]) + "' as prefix"
		}
	}

	// When a valid word or we have nowhere to go
	if wb.IsValidWord || !foundValidContinuation {
		// If we found no valid continuations, we're at a dead end
		// DO NOT add fallback letters - leave the sets empty
		wb.Suggestion = "This isn't a valid prefix or suffix of any word. Try removing some letters."
	}

}

// GetCurrentState returns the current state as a map
func (wb *EnhancedWordBuilder) GetCurrentState() map[string]interface{} {
	prefixSet := make([]string, 0, len(wb.PrefixSet))
	for letter := range wb.PrefixSet {
		prefixSet = append(prefixSet, letter)
	}

	suffixSet := make([]string, 0, len(wb.SuffixSet))
	for letter := range wb.SuffixSet {
		suffixSet = append(suffixSet, letter)
	}

	// Only return a few completions to avoid overwhelming the UI
	var displayCompletions []string
	if len(wb.ValidCompletions) > 5 {
		displayCompletions = wb.ValidCompletions[:5]
	} else {
		displayCompletions = wb.ValidCompletions
	}

	return map[string]interface{}{
		"answer":            wb.Answer,
		"prefix_set":        prefixSet,
		"suffix_set":        suffixSet,
		"step":              wb.Step,
		"is_valid_word":     wb.IsValidWord,
		"valid_completions": displayCompletions,
		"suggestion":        wb.Suggestion,
	}
}

// Helper function to determine minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
