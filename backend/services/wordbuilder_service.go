package services

import (
	"wordbuilder/models"
)

// WordBuilderService handles game state and operations
type WordBuilderService struct {
	Dictionary *models.WordDictionary
	Sessions   map[string]*models.EnhancedWordBuilder
}

// NewWordBuilderService creates a new service instance
func NewWordBuilderService(dictionary *models.WordDictionary) *WordBuilderService {
	return &WordBuilderService{
		Dictionary: dictionary,
		Sessions:   make(map[string]*models.EnhancedWordBuilder),
	}
}

// CreateSession initializes a new game session
func (s *WordBuilderService) CreateSession(sessionID string, dictService *DictionaryService) *models.EnhancedWordBuilder {
	// Safety check - ensure dictionary exists
	if s.Dictionary == nil {
		// Use provided dictionary service for fallback
		wordList, _ := dictService.LoadWordList("words.txt")
		s.Dictionary = dictService.CreateDictionary(wordList)
	}

	builder := models.NewEnhancedWordBuilder(s.Dictionary)
	s.Sessions[sessionID] = builder
	return builder
}

// GetSession retrieves a session by ID
func (s *WordBuilderService) GetSession(sessionID string) (*models.EnhancedWordBuilder, bool) {
	builder, exists := s.Sessions[sessionID]
	return builder, exists
}

// ResetSession resets a specific game session
func (s *WordBuilderService) ResetSession(sessionID string) (*models.EnhancedWordBuilder, bool) {
	builder, exists := s.Sessions[sessionID]
	if !exists {
		return nil, false
	}

	builder.Reset()
	return builder, true
}

// UpdateDictionary updates the dictionary used by the service and resets all active sessions
func (s *WordBuilderService) UpdateDictionary(dictionary *models.WordDictionary) {
	s.Dictionary = dictionary

	// Update the dictionary for all active sessions and reset them
	for _, builder := range s.Sessions {
		builder.Dictionary = dictionary
		builder.Reset() // Reset the state to work with the new dictionary
	}
}
