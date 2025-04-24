package models

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"

	utils "wordbuilder/utils"
)

// WordBuilderState holds the state for the word builder game, without any dependencies or methods.
type WordBuilderState struct {
	Answer           string
	PrefixSet        map[string]bool
	SuffixSet        map[string]bool
	Step             int
	IsValidWord      bool
	ValidCompletions []string
	Suggestion       string
}

type TrieI interface {
	GetNextLetters(prefix string) []string
	KeysWithPrefix(prefix string) []string
}

// WordDictionary defines the methods required by the word builder logic.
type WordDictionaryI interface {
	ContainsWord(word string) bool
	FindWordsWithPrefix(prefix string) []string
	FindWordsWithSuffix(suffix string) []string
	GetForwardTrie() TrieI
	GetReverseTrie() TrieI
	GetWordList() []string
}

// CheckValidWord verifies if the current answer is a valid word
func CheckValidWord(state WordBuilderState, dict WordDictionaryI) bool {
	return len(state.Answer) > 0 && dict.ContainsWord(state.Answer)
}

// AddLetter adds a letter to either the prefix or suffix
func AddLetter(state WordBuilderState, dict WordDictionaryI, letter, position string) (WordBuilderState, string, error) {
	// Work with a copy of state
	newState := state
	letter = strings.ToLower(letter)

	if position == "prefix" {
		if !state.PrefixSet[letter] {
			return state, fmt.Sprintf("Invalid letter '%s' for prefix position.", letter), fmt.Errorf("invalid prefix letter")
		}
		newState.Answer = letter + state.Answer
	} else if position == "suffix" {
		if !state.SuffixSet[letter] {
			return state, fmt.Sprintf("Invalid letter '%s' for suffix position.", letter), fmt.Errorf("invalid suffix letter")
		}
		newState.Answer = state.Answer + letter
	} else {
		return state, "Invalid position. Use 'prefix' or 'suffix'.", fmt.Errorf("invalid position")
	}

	// Call pure versions of CheckValidWord and UpdateSets
	newState.IsValidWord = CheckValidWord(newState, dict)
	newState = UpdateSets(newState, dict)
	newState.Step++

	message := fmt.Sprintf("Step %d: Added '%s' as %s -> Answer: %s", newState.Step, letter, position, newState.Answer)
	if newState.IsValidWord {
		message += fmt.Sprintf("\n*** '%s' is a valid word! ***", newState.Answer)
	}

	// Add information about possible completions
	if len(newState.ValidCompletions) > 0 && !newState.IsValidWord {
		message += fmt.Sprintf("\nPossible completions: %s", strings.Join(newState.ValidCompletions[:min(3, len(newState.ValidCompletions))], ", "))
	}

	return newState, message, nil
}

// RemoveLetter removes a letter at the specified index
func RemoveLetter(state WordBuilderState, dict WordDictionaryI, index int) (WordBuilderState, string, error) {
	if index < 0 || index >= len(state.Answer) {
		return state, fmt.Sprintf("Invalid index %d for answer '%s'.", index, state.Answer), fmt.Errorf("invalid index")
	}

	newState := state
	letter := string(state.Answer[index])
	newState.Answer = state.Answer[:index] + state.Answer[index+1:]

	newState.IsValidWord = CheckValidWord(newState, dict)
	newState = UpdateSets(newState, dict)
	newState.Step++

	message := fmt.Sprintf("Step %d: Removed '%s' at index %d -> Answer: %s", newState.Step, letter, index, newState.Answer)
	if newState.IsValidWord {
		message += fmt.Sprintf("\n*** '%s' is a valid word! ***", newState.Answer)
	}

	return newState, message, nil
}

// UpdateSets updates the available prefix and suffix letter sets
func UpdateSets(state WordBuilderState, dict WordDictionaryI) WordBuilderState {
	newState := state
	newState.PrefixSet = make(map[string]bool)
	newState.SuffixSet = make(map[string]bool)
	newState.ValidCompletions = []string{}

	// If no letters yet, provide all letters that can start or end words
	if len(newState.Answer) == 0 {
		for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
			letterStr := string(letter)
			if len(dict.FindWordsWithPrefix(letterStr)) > 0 {
				newState.PrefixSet[letterStr] = true
			}
			if len(dict.FindWordsWithSuffix(letterStr)) > 0 {
				newState.SuffixSet[letterStr] = true
			}
		}
		return newState
	}

	foundValidContinuation := false

	// 1. Suffix letters from ForwardTrie
	suffixLetters := dict.GetForwardTrie().GetNextLetters(newState.Answer)
	for _, letter := range suffixLetters {
		newState.SuffixSet[letter] = true
		foundValidContinuation = true
	}
	if len(suffixLetters) > 0 {
		words := dict.FindWordsWithPrefix(newState.Answer)
		for i, word := range words {
			if i >= 5 {
				break
			}
			if len(word) > len(newState.Answer) {
				newState.ValidCompletions = append(newState.ValidCompletions, word)
			}
		}
	}

	// 2. Prefix letters from ReverseTrie
	// Assuming you have a ReverseString utility function
	reversedAnswer := utils.ReverseString(newState.Answer)
	prefixLetters := dict.GetReverseTrie().GetNextLetters(reversedAnswer)
	for _, letter := range prefixLetters {
		newState.PrefixSet[letter] = true
		foundValidContinuation = true
	}
	if len(prefixLetters) > 0 {
		revWords := dict.GetReverseTrie().KeysWithPrefix(reversedAnswer)
		for i, revWord := range revWords {
			if i >= 5 {
				break
			}
			if len(revWord) > len(reversedAnswer) {
				newState.ValidCompletions = append(newState.ValidCompletions, utils.ReverseString(revWord))
			}
		}
	}

	// 3. Embedded check with parallelism
	wordList := dict.GetWordList()
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
				idx := strings.Index(word, newState.Answer)
				if idx >= 0 {
					foundValidContinuation = true
					if idx > 0 {
						localPrefix[string(word[idx-1])] = true
					}
					embedEndIdx := idx + len(newState.Answer)
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
			newState.PrefixSet[letter] = true
		}
		for letter := range res.suffixSet {
			newState.SuffixSet[letter] = true
		}
	}

	// Generate suggestion (simplified)
	if len(newState.ValidCompletions) > 0 && !newState.IsValidWord {
		sort.Slice(newState.ValidCompletions, func(i, j int) bool {
			return len(newState.ValidCompletions[i]) < len(newState.ValidCompletions[j])
		})
		shortest := newState.ValidCompletions[0]
		if strings.HasPrefix(shortest, newState.Answer) && len(shortest) > len(newState.Answer) {
			newState.Suggestion = "Try adding '" + string(shortest[len(newState.Answer)]) + "' as suffix"
		} else if idx := strings.Index(shortest, newState.Answer); idx > 0 {
			newState.Suggestion = "Try adding '" + string(shortest[idx-1]) + "' as prefix"
		}
	}

	// When a valid word or we have nowhere to go
	if newState.IsValidWord || !foundValidContinuation {
		newState.Suggestion = "This isn't a valid prefix or suffix of any word. Try removing some letters."
	}

	return newState
}

// GetCurrentState returns the current state as a map
func GetCurrentState(state WordBuilderState) map[string]interface{} {
	prefixSet := make([]string, 0, len(state.PrefixSet))
	for letter := range state.PrefixSet {
		prefixSet = append(prefixSet, letter)
	}

	suffixSet := make([]string, 0, len(state.SuffixSet))
	for letter := range state.SuffixSet {
		suffixSet = append(suffixSet, letter)
	}

	// Only return a few completions to avoid overwhelming the UI
	var displayCompletions []string
	if len(state.ValidCompletions) > 5 {
		displayCompletions = state.ValidCompletions[:5]
	} else {
		displayCompletions = state.ValidCompletions
	}

	return map[string]interface{}{
		"answer":            state.Answer,
		"prefix_set":        prefixSet,
		"suffix_set":        suffixSet,
		"step":              state.Step,
		"is_valid_word":     state.IsValidWord,
		"valid_completions": displayCompletions,
		"suggestion":        state.Suggestion,
	}
}
