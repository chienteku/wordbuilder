package services

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	utils "wordbuilder/utils"
)

func TestLoadWordList(t *testing.T) {
	// Create a new dictionary service for testing
	service := NewDictionaryService()

	t.Run("successful load with mixed case and whitespace", func(t *testing.T) {
		// Create a temporary file with test words
		tempDir := t.TempDir()
		filename := filepath.Join(tempDir, "test_dict.txt")

		// Write test content to file
		content := []byte("Apple\nBanana \n CHERRY\n\ndate\n  Elderberry  ")
		err := os.WriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Expected result after processing
		expected := []string{"apple", "banana", "cherry", "date", "elderberry"}

		// Call the function being tested
		words, err := service.LoadWordList(filename)
		if err != nil {
			t.Fatalf("LoadWordList failed: %v", err)
		}

		// Compare results
		if !reflect.DeepEqual(words, expected) {
			t.Errorf("LoadWordList returned %v, want %v", words, expected)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		// Test with non-existing file
		_, err := service.LoadWordList("nonexistent_file.txt")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		// Create empty file
		tempDir := t.TempDir()
		filename := filepath.Join(tempDir, "empty.txt")

		err := os.WriteFile(filename, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create empty test file: %v", err)
		}

		// Test with empty file
		words, err := service.LoadWordList(filename)
		if err != nil {
			t.Fatalf("LoadWordList failed with empty file: %v", err)
		}

		if len(words) != 0 {
			t.Errorf("Expected empty result for empty file, got %v", words)
		}
	})

	t.Run("file with only whitespace", func(t *testing.T) {
		// Create file with only whitespace
		tempDir := t.TempDir()
		filename := filepath.Join(tempDir, "whitespace.txt")

		err := os.WriteFile(filename, []byte("  \n\t\n  \n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create whitespace test file: %v", err)
		}

		// Test with whitespace file
		words, err := service.LoadWordList(filename)
		if err != nil {
			t.Fatalf("LoadWordList failed with whitespace file: %v", err)
		}

		if len(words) != 0 {
			t.Errorf("Expected empty result for whitespace-only file, got %v", words)
		}
	})

	t.Run("large file", func(t *testing.T) {
		// Create a larger file to test scanning performance
		tempDir := t.TempDir()
		filename := filepath.Join(tempDir, "large.txt")

		// Create a file with 1000 words
		file, err := os.Create(filename)
		if err != nil {
			t.Fatalf("Failed to create large test file: %v", err)
		}

		expectedWords := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			word := "word" + string(rune('a'+i%26))
			file.WriteString(word + "\n")
			expectedWords[i] = word
		}
		file.Close()

		// Test with large file
		words, err := service.LoadWordList(filename)
		if err != nil {
			t.Fatalf("LoadWordList failed with large file: %v", err)
		}

		if len(words) != 1000 {
			t.Errorf("Expected 1000 words, got %d", len(words))
		}
	})

	t.Run("file with special characters", func(t *testing.T) {
		// Create file with special characters
		tempDir := t.TempDir()
		filename := filepath.Join(tempDir, "special.txt")

		content := []byte("Café\nnaïve\nco-operate\n")
		err := os.WriteFile(filename, content, 0644)
		if err != nil {
			t.Fatalf("Failed to create special chars test file: %v", err)
		}

		expected := []string{"café", "naïve", "co-operate"}

		// Test with special characters
		words, err := service.LoadWordList(filename)
		if err != nil {
			t.Fatalf("LoadWordList failed with special chars: %v", err)
		}

		if !reflect.DeepEqual(words, expected) {
			t.Errorf("Special chars test: got %v, want %v", words, expected)
		}
	})
}

func TestCreateDictionary(t *testing.T) {
	// Create a new dictionary service for testing
	service := NewDictionaryService()

	t.Run("empty word list", func(t *testing.T) {
		// Test with an empty word list
		wordList := []string{}
		dict := service.CreateDictionary(wordList)

		// Verify dictionary structure
		if dict == nil {
			t.Fatal("CreateDictionary returned nil for empty word list")
		}

		// Check that all structures are initialized but empty
		if len(dict.WordSet) != 0 {
			t.Errorf("Expected empty WordSet, got size %d", len(dict.WordSet))
		}

		if len(dict.WordList) != 0 {
			t.Errorf("Expected empty WordList, got size %d", len(dict.WordList))
		}

		// Verify that tries are initialized
		if dict.ForwardTrie == nil {
			t.Error("ForwardTrie was not initialized")
		}

		if dict.ReverseTrie == nil {
			t.Error("ReverseTrie was not initialized")
		}
	})

	t.Run("word list with mixed case", func(t *testing.T) {
		// Test with a word list containing mixed case
		wordList := []string{"Apple", "banana", "CHERRY"}
		dict := service.CreateDictionary(wordList)

		// Check that words were converted to lowercase in WordSet
		expectedWordSet := map[string]bool{
			"apple":  true,
			"banana": true,
			"cherry": true,
		}

		if !reflect.DeepEqual(dict.WordSet, expectedWordSet) {
			t.Errorf("WordSet does not match expected. Got %v, want %v", dict.WordSet, expectedWordSet)
		}

		// Check that WordList contains all lowercase words
		if len(dict.WordList) != 3 {
			t.Errorf("Expected WordList size 3, got %d", len(dict.WordList))
		}

		// Since map iteration order is not guaranteed, we'll check that all words exist
		for _, word := range []string{"apple", "banana", "cherry"} {
			found := false
			for _, dictWord := range dict.WordList {
				if dictWord == word {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Word '%s' not found in WordList", word)
			}
		}
	})

	t.Run("trie population", func(t *testing.T) {
		// Test that tries are populated correctly
		wordList := []string{"cat", "dog"}
		dict := service.CreateDictionary(wordList)

		// Test forward trie
		if !dict.ForwardTrie.Contains("cat") {
			t.Error("Word 'cat' not found in ForwardTrie")
		}
		if !dict.ForwardTrie.Contains("dog") {
			t.Error("Word 'dog' not found in ForwardTrie")
		}
		if dict.ForwardTrie.Contains("bat") {
			t.Error("Word 'bat' incorrectly found in ForwardTrie")
		}

		// Test reverse trie (words are reversed in the ReverseTrie)
		if !dict.ReverseTrie.Contains("tac") { // Reverse of "cat"
			t.Error("Reversed word 'tac' not found in ReverseTrie")
		}
		if !dict.ReverseTrie.Contains("god") { // Reverse of "dog"
			t.Error("Reversed word 'god' not found in ReverseTrie")
		}
		if dict.ReverseTrie.Contains("tab") { // Reverse of "bat" (not in dictionary)
			t.Error("Reversed word 'tab' incorrectly found in ReverseTrie")
		}
	})

	t.Run("prefix search in tries", func(t *testing.T) {
		// Test prefix search functionality of tries
		wordList := []string{"cat", "car", "cart", "dog"}
		dict := service.CreateDictionary(wordList)

		// Test forward prefix
		if len(dict.ForwardTrie.KeysWithPrefix("ca")) == 0 {
			t.Error("Prefix 'ca' not found in ForwardTrie")
		}
		if len(dict.ForwardTrie.KeysWithPrefix("ba")) != 0 {
			t.Error("Prefix 'ba' incorrectly found in ForwardTrie")
		}

		// Test reverse prefix (remember words are reversed in ReverseTrie)
		if len(dict.ReverseTrie.KeysWithPrefix("ra")) == 0 { // Reverse prefix of "ar" ending
			t.Error("Reversed prefix 'ra' not found in ReverseTrie")
		}
		if len(dict.ReverseTrie.KeysWithPrefix("la")) != 0 {
			t.Error("Reversed prefix 'la' incorrectly found in ReverseTrie")
		}
	})

	t.Run("duplicate words", func(t *testing.T) {
		// Test handling of duplicate words
		wordList := []string{"cat", "dog", "cat", "cat"}
		dict := service.CreateDictionary(wordList)

		// WordSet should deduplicate
		if len(dict.WordSet) != 2 {
			t.Errorf("Expected WordSet size 2 after deduplication, got %d", len(dict.WordSet))
		}

		// WordList should keep duplicates
		if len(dict.WordList) != 4 {
			t.Errorf("Expected WordList size 4 with duplicates, got %d", len(dict.WordList))
		}

		// Count occurrences of "cat" in WordList
		catCount := 0
		for _, word := range dict.WordList {
			if word == "cat" {
				catCount++
			}
		}
		if catCount != 3 {
			t.Errorf("Expected 3 occurrences of 'cat' in WordList, got %d", catCount)
		}
	})

	t.Run("special characters", func(t *testing.T) {
		// Test with words containing special characters
		wordList := []string{"café", "naïve", "co-op"}
		dict := service.CreateDictionary(wordList)

		// Check WordSet contains special characters
		for _, word := range wordList {
			if !dict.WordSet[word] {
				t.Errorf("Word '%s' not found in WordSet", word)
			}
		}

		// Test forward trie with special characters
		if !dict.ForwardTrie.Contains("café") {
			t.Error("Word 'café' not found in ForwardTrie")
		}

		// Test reverse trie with special characters
		reversedCafe := utils.ReverseString("café")
		if !dict.ReverseTrie.Contains(reversedCafe) {
			t.Errorf("Reversed word '%s' not found in ReverseTrie", reversedCafe)
		}
	})

	t.Run("large dictionary", func(t *testing.T) {
		// Test with a larger word list
		wordList := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			wordList[i] = "word" + string(rune('a'+i%26))
		}

		dict := service.CreateDictionary(wordList)

		// Basic size checks
		if len(dict.WordList) != 1000 {
			t.Errorf("Expected WordList size 1000, got %d", len(dict.WordList))
		}

		// Since there are only 26 unique words (worda through wordz),
		// the WordSet should only have 26 entries
		if len(dict.WordSet) != 26 {
			t.Errorf("Expected WordSet size 26 after deduplication, got %d", len(dict.WordSet))
		}

		// Check that all expected words are in the forward trie
		for i := 0; i < 26; i++ {
			word := "word" + string(rune('a'+i))
			if !dict.ForwardTrie.Contains(word) {
				t.Errorf("Word '%s' not found in ForwardTrie", word)
			}
		}

		// Check a word that shouldn't be in the trie
		if dict.ForwardTrie.Contains("wordaa") {
			t.Error("Word 'wordaa' incorrectly found in ForwardTrie")
		}
	})
}
