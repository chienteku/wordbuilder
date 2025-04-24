package main

import (
	"testing"
	. "wordbuilder/models"
	"wordbuilder/services"
)

// Benchmark for dictionary creation
func BenchmarkDictionaryCreation(b *testing.B) {
	dictService := services.NewDictionaryService()
	wordList, _ := dictService.LoadWordList("words.txt")

	b.ResetTimer() // Ignore setup time
	for i := 0; i < b.N; i++ {
		NewWordDictionary(wordList)
	}
}

// Benchmark for initial WordBuilder state
func BenchmarkInitialWordBuilderState(b *testing.B) {
	dictService := services.NewDictionaryService()
	wordList, _ := dictService.LoadWordList("words.txt")
	dictionary := NewWordDictionary(wordList)

	b.ResetTimer() // Ignore setup time
	for i := 0; i < b.N; i++ {
		state := WordBuilderState{
			Answer:           "",
			PrefixSet:        make(map[string]bool),
			SuffixSet:        make(map[string]bool),
			Step:             0,
			IsValidWord:      false,
			ValidCompletions: []string{},
			Suggestion:       "",
		}
		_ = UpdateSets(state, dictionary)
	}
}

// Benchmark for single-letter suffix searches
func BenchmarkSingleLetterSuffixSearch(b *testing.B) {
	dictService := services.NewDictionaryService()
	wordList, _ := dictService.LoadWordList("words.txt")
	dictionary := NewWordDictionary(wordList)

	// Test different letters
	for _, letter := range []string{"a", "e", "s", "z"} {
		b.Run("Letter_"+letter, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dictionary.FindWordsWithSuffix(letter)
			}
		})
	}
}
