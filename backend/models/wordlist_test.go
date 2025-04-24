package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestWordListInitialization(t *testing.T) {
	now := time.Now()
	wl := WordList{
		ID:          1,
		Name:        "Sample List",
		Description: "A test word list",
		Source:      "manual",
		FilePath:    "/tmp/words.txt",
		WordCount:   10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if wl.ID != 1 || wl.Name != "Sample List" || wl.WordCount != 10 {
		t.Errorf("WordList fields not set correctly: %+v", wl)
	}
}

func TestWordListJSONMarshaling(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	wl := WordList{
		ID:          2,
		Name:        "Another List",
		Description: "JSON test",
		Source:      "imported",
		FilePath:    "/tmp/another.txt",
		WordCount:   20,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, err := json.Marshal(wl)
	if err != nil {
		t.Fatalf("Failed to marshal WordList: %v", err)
	}

	var wl2 WordList
	if err := json.Unmarshal(data, &wl2); err != nil {
		t.Fatalf("Failed to unmarshal WordList: %v", err)
	}

	// Compare all fields except time.Time (which may lose precision in JSON)
	if wl.ID != wl2.ID || wl.Name != wl2.Name || wl.Description != wl2.Description ||
		wl.Source != wl2.Source || wl.FilePath != wl2.FilePath || wl.WordCount != wl2.WordCount {
		t.Errorf("WordList JSON roundtrip mismatch: got %+v, want %+v", wl2, wl)
	}
	// For time fields, allow for second-level precision
	if !wl.CreatedAt.Equal(wl2.CreatedAt) || !wl.UpdatedAt.Equal(wl2.UpdatedAt) {
		t.Errorf("WordList time fields mismatch: got %+v, want %+v", wl2, wl)
	}
}
