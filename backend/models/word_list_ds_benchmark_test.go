package models

import (
	"math/rand"
	"testing"
)

// generateWordList generates a list of n random words with lengths between 5 and 20 characters.
func generateWordList(n int) []string {
	rand.Seed(42) // For reproducibility
	var words []string
	for i := 0; i < n; i++ {
		len := 5 + rand.Intn(16) // Length between 5 and 20
		word := make([]byte, len)
		for j := 0; j < len; j++ {
			word[j] = 'a' + byte(rand.Intn(26))
		}
		words = append(words, string(word))
	}
	return words
}

// BenchmarkNewWordListDS measures the time and allocations to create a new WordListDS with 100,000 words.
func BenchmarkNewWordListDS(b *testing.B) {
	words := generateWordList(100000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewWordListDS(words)
	}
}

// BenchmarkGetGroups measures the performance of GetGroups for different query lengths.
func BenchmarkGetGroups(b *testing.B) {
	words := generateWordList(100000)
	ds := NewWordListDS(words)
	b.Run("empty", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ds.GetGroups("")
		}
	})
	b.Run("single", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ds.GetGroups("a")
		}
	})
	b.Run("multi", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ds.GetGroups("ab")
		}
	})
	b.Run("long", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = ds.GetGroups("abcde")
		}
	})
}

// BenchmarkContains measures the performance of Contains for existing and non-existing words.
func BenchmarkContains(b *testing.B) {
	words := generateWordList(100000)
	ds := NewWordListDS(words)
	existing := words[0]
	nonExisting := "nonexistingword" // Assumed to be non-existing
	b.Run("existing", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ds.Contains(existing)
		}
	})
	b.Run("non-existing", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ds.Contains(nonExisting)
		}
	})
}

// BenchmarkGetSuggestionGroups measures the performance of GetSuggestionGroups for different query lengths.
func BenchmarkGetSuggestionGroups(b *testing.B) {
	words := generateWordList(100000)
	ds := NewWordListDS(words)
	b.Run("empty", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = ds.GetSuggestionGroups("")
		}
	})
	b.Run("single", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = ds.GetSuggestionGroups("a")
		}
	})
	b.Run("multi", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = ds.GetSuggestionGroups("ab")
		}
	})
	b.Run("long", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = ds.GetSuggestionGroups("abcde")
		}
	})
}
