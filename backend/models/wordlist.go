package models

import (
	"time"
)

// WordList represents a list of words with metadata
type WordList struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	FilePath    string    `json:"file_path"`
	WordCount   int       `json:"word_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
