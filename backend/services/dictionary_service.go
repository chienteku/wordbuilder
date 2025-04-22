package services

import (
	"bufio"
	"os"
	"strings"

	"wordbuilder/models"
	utils "wordbuilder/utils"
)

// DictionaryService provides operations for dictionary functionality
type DictionaryService struct {
	// Any necessary fields
}

// NewDictionaryService creates a new dictionary service
func NewDictionaryService() *DictionaryService {
	return &DictionaryService{}
}

// LoadWordList loads the dictionary from a file
func (s *DictionaryService) LoadWordList(filename string) ([]string, error) {
	// Move this function from main.go
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if word != "" {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

// CreateDictionary creates a new WordDictionary from a word list
func (s *DictionaryService) CreateDictionary(wordList []string) *models.WordDictionary {
	dict := &models.WordDictionary{
		WordSet:     make(map[string]bool),
		ForwardTrie: models.NewTrie(),
		ReverseTrie: models.NewTrie(),
		WordList:    make([]string, 0, len(wordList)),
	}

	for _, word := range wordList {
		word = strings.ToLower(word)
		dict.WordSet[word] = true
		dict.ForwardTrie.Insert(word)
		dict.ReverseTrie.Insert(utils.ReverseString(word))
		dict.WordList = append(dict.WordList, word)
	}

	return dict
}
