package models

import (
	"reflect"
	"testing"
)

// TestNewWordListDS tests the creation of a new WordListDS.
func TestNewWordListDS(t *testing.T) {
	words := []string{"elephant", "envelope", "pen", "penguin", "people", "person", "personal", "prepositions", "repeat", "sheep", "sleep"}
	ds := NewWordListDS(words)

	// Check if data is correctly formatted
	expectedData := []byte("elephant$envelope$pen$penguin$people$person$personal$prepositions$repeat$sheep$sleep")
	if !reflect.DeepEqual(ds.data, expectedData) {
		t.Errorf("Expected data %v, got %v", expectedData, ds.data)
	}

	// Check if index is not nil
	if ds.index == nil {
		t.Error("Expected non-nil suffixarray index")
	}
}

// TestGetGroups_EmptyQuery tests GetGroups with an empty query string.
func TestGetGroups_EmptyQuery(t *testing.T) {
	words := []string{"elephant", "envelope", "pen", "penguin", "people", "person", "personal", "prepositions", "repeat", "sheep", "sleep", "zoo"}
	ds := NewWordListDS(words)

	pre, app := ds.GetGroups("")

	expectedPre := []string{"a", "e", "g", "h", "i", "l", "n", "o", "p", "r", "s", "t", "u", "v", "z"}
	expectedApp := []string{"a", "e", "g", "h", "i", "l", "n", "o", "p", "r", "s", "t", "u", "v"}

	if !reflect.DeepEqual(pre, expectedPre) {
		t.Errorf("Expected prepend %v, got %v", expectedPre, pre)
	}
	if !reflect.DeepEqual(app, expectedApp) {
		t.Errorf("Expected append %v, got %v", expectedApp, app)
	}
}

// TestGetGroups_SingleCharacterQuery tests GetGroups with a single character query.
func TestGetGroups_SingleCharacterQuery(t *testing.T) {
	words := []string{"elephant", "envelope", "pen", "penguin", "people", "person", "personal", "prepositions", "repeat", "sheep", "sleep"}
	ds := NewWordListDS(words)

	pre, app := ds.GetGroups("e")
	expectedPre := []string{"e", "h", "l", "p", "r", "v"}
	expectedApp := []string{"a", "e", "l", "n", "o", "p", "r"}

	if !reflect.DeepEqual(pre, expectedPre) {
		t.Errorf("Expected prepend %v, got %v", expectedPre, pre)
	}
	if !reflect.DeepEqual(app, expectedApp) {
		t.Errorf("Expected append %v, got %v", expectedApp, app)
	}
}

// TestGetGroups_MultiCharacterQuery tests GetGroups with a multi-character query.
func TestGetGroups_MultiCharacterQuery(t *testing.T) {
	words := []string{"elephant", "envelope", "pen", "penguin", "people", "person", "personal", "prepositions", "repeat", "sheep", "sleep"}
	ds := NewWordListDS(words)

	pre, app := ds.GetGroups("ep")
	expectedPre := []string{"e", "l", "r"}
	expectedApp := []string{"e", "h", "o"}

	if !reflect.DeepEqual(pre, expectedPre) {
		t.Errorf("Expected prepend %v, got %v", expectedPre, pre)
	}
	if !reflect.DeepEqual(app, expectedApp) {
		t.Errorf("Expected append %v, got %v", expectedApp, app)
	}
}

// TestGetGroups_NoMatch tests GetGroups with a query that doesn't exist in any word.
func TestGetGroups_NoMatch(t *testing.T) {
	words := []string{"elephant", "envelope", "pen", "penguin", "people", "person", "personal", "prepositions", "repeat", "sheep", "sleep"}
	ds := NewWordListDS(words)

	pre, app := ds.GetGroups("xyz")
	if len(pre) != 0 || len(app) != 0 {
		t.Errorf("Expected empty prepend and append lists, got prepend=%v, append=%v", pre, app)
	}
}

// TestGetGroups_EntireWord tests GetGroups with a query that matches an entire word.
func TestGetGroups_EntireWord(t *testing.T) {
	words := []string{"elephant", "envelope", "pen", "penguin", "people", "person", "personal", "prepositions", "repeat", "sheep", "sleep"}
	ds := NewWordListDS(words)

	pre, app := ds.GetGroups("pen")
	expectedPre := []string{}
	expectedApp := []string{"g"}

	if !equalStringSlices(pre, expectedPre) {
		t.Errorf("Expected prepend %v, got %v", expectedPre, pre)
	}
	if !equalStringSlices(app, expectedApp) {
		t.Errorf("Expected append %v, got %v", expectedApp, app)
	}
}

func TestWordListDS_Contains(t *testing.T) {
	// Initialize a WordListDS with a sample word list
	words := []string{"elephant", "envelope", "pen", "penguin", "people", "person", "personal", "prepositions", "repeat", "sheep", "sleep"}
	ds := NewWordListDS(words)

	tests := []struct {
		name     string
		word     string
		expected bool
	}{
		{
			name:     "Existing word - middle",
			word:     "people",
			expected: true,
		},
		{
			name:     "Existing word - start",
			word:     "elephant",
			expected: true,
		},
		{
			name:     "Existing word - end",
			word:     "sleep",
			expected: true,
		},
		{
			name:     "Non-existing word",
			word:     "grape",
			expected: false,
		},
		{
			name:     "Empty string",
			word:     "",
			expected: false,
		},
		{
			name:     "Partial word match",
			word:     "ban",
			expected: false,
		},
		{
			name:     "Case sensitivity",
			word:     "Banana",
			expected: false,
		},
		{
			name:     "Word with suffix",
			word:     "bananas",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ds.Contains(tt.word)
			if result != tt.expected {
				t.Errorf("Contains(%q) = %v; want %v", tt.word, result, tt.expected)
			}
		})
	}

	// Test with an empty WordListDS
	emptyDS := NewWordListDS([]string{})
	t.Run("Empty WordListDS", func(t *testing.T) {
		result := emptyDS.Contains("test")
		if result != false {
			t.Errorf("Contains(%q) on empty WordListDS = %v; want false", "test", result)
		}
	})

	// Test with a single-word WordListDS
	singleDS := NewWordListDS([]string{"solo"})
	t.Run("Single word - existing", func(t *testing.T) {
		result := singleDS.Contains("solo")
		if result != true {
			t.Errorf("Contains(%q) on single-word WordListDS = %v; want true", "solo", result)
		}
	})
	t.Run("Single word - non-existing", func(t *testing.T) {
		result := singleDS.Contains("other")
		if result != false {
			t.Errorf("Contains(%q) on single-word WordListDS = %v; want false", "other", result)
		}
	})
}
