package services

import (
	"wordbuilder/models"
)

// WordBuilderService handles game state and operations
type WordBuilderService struct {
	Dictionary *models.WordDictionary
	Sessions   map[string]*models.WordBuilderState
	// or models.WordBuilderState if you don't want pointers
}

// NewWordBuilderService creates a new service instance
func NewWordBuilderService(dictionary *models.WordDictionary) *WordBuilderService {
	return &WordBuilderService{
		Dictionary: dictionary,
		Sessions:   make(map[string]*models.WordBuilderState),
	}
}

// CreateSession initializes a new game session
func (s *WordBuilderService) CreateSession(sessionID string, dictService *DictionaryService) *models.WordBuilderState {
	// Safety check - ensure dictionary exists
	if s.Dictionary == nil {
		wordList, _ := dictService.LoadWordList("words.txt")
		s.Dictionary = dictService.CreateDictionary(wordList)
	}

	// Initialize the state
	state := models.WordBuilderState{
		Answer:           "",
		PrefixSet:        make(map[string]bool),
		SuffixSet:        make(map[string]bool),
		Step:             0,
		IsValidWord:      false,
		ValidCompletions: []string{},
		Suggestion:       "",
	}
	// Initialize sets using the pure function
	state = models.UpdateSets(state, s.Dictionary)
	s.Sessions[sessionID] = &state
	return &state
}

// GetSession retrieves a session by ID
func (s *WordBuilderService) GetSession(sessionID string) (*models.WordBuilderState, bool) {
	builder, exists := s.Sessions[sessionID]
	return builder, exists
}

// ResetSession resets a specific game session
func (s *WordBuilderService) ResetSession(sessionID string) (*models.WordBuilderState, bool) {
	_, exists := s.Sessions[sessionID]
	if !exists {
		return nil, false
	}
	// Create a new, fresh state
	state := models.WordBuilderState{
		Answer:           "",
		PrefixSet:        make(map[string]bool),
		SuffixSet:        make(map[string]bool),
		Step:             0,
		IsValidWord:      false,
		ValidCompletions: []string{},
		Suggestion:       "",
	}

	state = models.UpdateSets(state, s.Dictionary)
	s.Sessions[sessionID] = &state
	return &state, true
}

// UpdateDictionary updates the dictionary used by the service and resets all active sessions
func (s *WordBuilderService) UpdateDictionary(dictionary *models.WordDictionary) {
	s.Dictionary = dictionary

	// Reset all active sessions to use the new dictionary
	for sessionID := range s.Sessions {
		state := models.WordBuilderState{
			Answer:           "",
			PrefixSet:        make(map[string]bool),
			SuffixSet:        make(map[string]bool),
			Step:             0,
			IsValidWord:      false,
			ValidCompletions: []string{},
			Suggestion:       "",
		}
		state = models.UpdateSets(state, s.Dictionary)
		s.Sessions[sessionID] = &state
	}
}
