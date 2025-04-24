package models

import (
	"testing"
	"time"
)

func TestWordListStruct(t *testing.T) {
	now := time.Now()
	wl := WordList{
		ID:          1,
		Name:        "Test",
		Description: "desc",
		Source:      "src",
		FilePath:    "/tmp/file",
		WordCount:   10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if wl.ID != 1 || wl.Name != "Test" || wl.WordCount != 10 {
		t.Error("WordList struct fields not set correctly")
	}
}
