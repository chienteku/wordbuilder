package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wordbuilder/models"

	lru "github.com/hashicorp/golang-lru"
)

// WordListService handles operations for word lists
type WordListService struct {
	DBService         *DatabaseService
	DictionaryService *DictionaryService
	UploadDir         string
	dictCache         *lru.Cache
}

// NewWordListService creates a new word list service
func NewWordListService(dbService *DatabaseService, dictService *DictionaryService, uploadDir string) *WordListService {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	// Initialize LRU cache
	cache, err := lru.New(100) // Adjust size as needed

	if err != nil {
		panic(fmt.Sprintf("Failed to create LRU cache: %v", err))
	}

	return &WordListService{
		DBService:         dbService,
		DictionaryService: dictService,
		UploadDir:         uploadDir,
		dictCache:         cache,
	}
}

// CreateWordList saves an uploaded word list file and metadata
func (s *WordListService) CreateWordList(fileData []byte, name, description, source string) (*models.WordList, error) {
	// Generate a unique filename
	timestamp := time.Now().UnixNano()
	sanitizedName := strings.ReplaceAll(name, " ", "_")
	filename := fmt.Sprintf("%s_%d.txt", sanitizedName, timestamp)
	filepath := filepath.Join(s.UploadDir, filename)

	// Write file to disk
	err := os.WriteFile(filepath, fileData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Count words and validate content
	wordCount, err := s.countWordsInFile(filepath)
	if err != nil {
		// Clean up file if error occurs
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to process word list: %w", err)
	}

	// Create word list entry
	wordList := &models.WordList{
		Name:        name,
		Description: description,
		Source:      source,
		FilePath:    filepath,
		WordCount:   wordCount,
	}

	// Insert into database
	id, err := s.DBService.InsertWordList(wordList)
	if err != nil {
		// Clean up file if error occurs
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to save word list metadata: %w", err)
	}

	wordList.ID = id
	return wordList, nil
}

// UpdateWordList updates a word list and optionally replaces the file
func (s *WordListService) UpdateWordList(id int, name, description, source string, fileData []byte) (*models.WordList, error) {
	// Get existing word list
	wordList, err := s.DBService.GetWordList(id)
	if err != nil {
		return nil, fmt.Errorf("word list not found: %w", err)
	}

	// Update metadata
	wordList.Name = name
	wordList.Description = description
	wordList.Source = source

	// If new file is provided, replace the existing one
	if fileData != nil && len(fileData) > 0 {
		// Remove old file
		oldFilePath := wordList.FilePath
		if _, err := os.Stat(oldFilePath); err == nil {
			if err := os.Remove(oldFilePath); err != nil {
				return nil, fmt.Errorf("failed to remove old file: %w", err)
			}
		}

		// Generate a new filename
		timestamp := time.Now().UnixNano()
		sanitizedName := strings.ReplaceAll(name, " ", "_")
		filename := fmt.Sprintf("%s_%d.txt", sanitizedName, timestamp)
		filepath := filepath.Join(s.UploadDir, filename)

		// Write new file
		if err := os.WriteFile(filepath, fileData, 0644); err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}

		// Update file path and word count
		wordCount, err := s.countWordsInFile(filepath)
		if err != nil {
			// Clean up file if error occurs
			os.Remove(filepath)
			return nil, fmt.Errorf("failed to process word list: %w", err)
		}

		wordList.FilePath = filepath
		wordList.WordCount = wordCount
	}

	// Update in database
	if err := s.DBService.UpdateWordList(wordList); err != nil {
		return nil, fmt.Errorf("failed to update word list: %w", err)
	}

	// Clear cache
	s.dictCache.Remove(wordList.ID)

	return wordList, nil
}

// DeleteWordList removes a word list and its file
func (s *WordListService) DeleteWordList(id int) error {
	// Get the word list to find the file path
	wordList, err := s.DBService.GetWordList(id)
	if err != nil {
		return fmt.Errorf("word list not found: %w", err)
	}

	// Remove the file
	if _, err := os.Stat(wordList.FilePath); err == nil {
		if err := os.Remove(wordList.FilePath); err != nil {
			return fmt.Errorf("failed to remove word list file: %w", err)
		}
	}

	// Delete from database
	if err := s.DBService.DeleteWordList(id); err != nil {
		return fmt.Errorf("failed to delete word list metadata: %w", err)
	}

	// Clear cache
	s.dictCache.Remove(id)

	return nil
}

// GetWordList retrieves a word list by ID
func (s *WordListService) GetWordList(id int) (*models.WordList, error) {
	return s.DBService.GetWordList(id)
}

// GetAllWordLists retrieves all word lists
func (s *WordListService) GetAllWordLists() ([]*models.WordList, error) {
	return s.DBService.GetAllWordLists()
}

// countWordsInFile counts the number of words in a file
func (s *WordListService) countWordsInFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

// LoadWordListIntoDictionary loads a word list into a dictionary
func (s *WordListService) LoadWordListIntoDictionary(wordListID int) (*models.WordDictionary, error) {

	// Check cache first
	if cached, ok := s.dictCache.Get(wordListID); ok {
		return cached.(*models.WordDictionary), nil
	}

	// Get the word list
	wordList, err := s.DBService.GetWordList(wordListID)
	if err != nil {
		return nil, fmt.Errorf("word list not found: %w", err)
	}

	// Load the word list
	words, err := s.DictionaryService.LoadWordList(wordList.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load word list: %w", err)
	}

	// Create dictionary
	dictionary := s.DictionaryService.CreateDictionary(words)

	// Add to cache
	s.dictCache.Add(wordListID, dictionary)

	return dictionary, nil
}

// ReadWordListContent reads the content of a word list file
func (s *WordListService) ReadWordListContent(id int, limit int) ([]string, error) {
	// Get the word list
	wordList, err := s.DBService.GetWordList(id)
	if err != nil {
		return nil, fmt.Errorf("word list not found: %w", err)
	}

	// Open the file
	file, err := os.Open(wordList.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open word list file: %w", err)
	}
	defer file.Close()

	// Read words up to the limit
	words := make([]string, 0, limit)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && len(words) < limit {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read word list: %w", err)
	}

	return words, nil
}

// ImportWordListFromReader imports words from a reader
func (s *WordListService) ImportWordListFromReader(reader io.Reader, name, description, source string) (*models.WordList, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp(s.UploadDir, "import-*.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Copy data to the file
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Generate a permanent filename
	timestamp := time.Now().UnixNano()
	sanitizedName := strings.ReplaceAll(name, " ", "_")
	filename := fmt.Sprintf("%s_%d.txt", sanitizedName, timestamp)
	filepath := filepath.Join(s.UploadDir, filename)

	// Close the temporary file before renaming
	tempFile.Close()

	// Rename the temporary file
	if err := os.Rename(tempFile.Name(), filepath); err != nil {
		os.Remove(tempFile.Name())
		return nil, fmt.Errorf("failed to save word list file: %w", err)
	}

	// Count words
	wordCount, err := s.countWordsInFile(filepath)
	if err != nil {
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to process word list: %w", err)
	}

	// Create word list entry
	wordList := &models.WordList{
		Name:        name,
		Description: description,
		Source:      source,
		FilePath:    filepath,
		WordCount:   wordCount,
	}

	// Insert into database
	id, err := s.DBService.InsertWordList(wordList)
	if err != nil {
		os.Remove(filepath)
		return nil, fmt.Errorf("failed to save word list metadata: %w", err)
	}

	wordList.ID = id
	return wordList, nil
}
