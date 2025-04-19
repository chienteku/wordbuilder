package main

import (
	"testing"
)

// Benchmark for dictionary creation
func BenchmarkDictionaryCreation(b *testing.B) {
	wordList, _ := loadWordList("words.txt")

	b.ResetTimer() // Ignore setup time
	for i := 0; i < b.N; i++ {
		NewWordDictionary(wordList)
	}
}

// Benchmark for initial WordBuilder state
func BenchmarkInitialWordBuilderState(b *testing.B) {
	wordList, _ := loadWordList("words.txt")
	dictionary := NewWordDictionary(wordList)

	b.ResetTimer() // Ignore setup time
	for i := 0; i < b.N; i++ {
		NewEnhancedWordBuilder(dictionary)
	}
}

// Benchmark for single-letter suffix searches
func BenchmarkSingleLetterSuffixSearch(b *testing.B) {
	wordList, _ := loadWordList("words.txt")
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
