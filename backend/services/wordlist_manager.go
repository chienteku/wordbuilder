package services

import (
	"os"
	"strings"
	"sync"
)

type WordListManager struct {
	words []string
	mu    sync.RWMutex
}

var (
	instance *WordListManager
	once     sync.Once
)

// GetWordListManager returns the singleton instance
func GetWordListManager() *WordListManager {
	once.Do(func() {
		instance = &WordListManager{}
	})
	return instance
}

// LoadFromFile replaces the word list with new data from a file
func (w *WordListManager) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	words := strings.Split(string(data), "\n")
	w.mu.Lock()
	w.words = words
	w.mu.Unlock()
	return nil
}

// GetWords returns a copy of the word list
func (w *WordListManager) GetWords() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	wordsCopy := make([]string, len(w.words))
	copy(wordsCopy, w.words)
	return wordsCopy
}
