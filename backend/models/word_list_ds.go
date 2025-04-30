package models

import (
	"bytes"
	"index/suffixarray"
	"sort"
	"strings"
)

// WordListDS is a data structure to store and query a large word list.
type WordListDS struct {
	data          []byte
	index         *suffixarray.Index
	initalPrepend []string
	initialAppend []string
}

// NewWordListDS creates a new WordListDS with the given word list.
func NewWordListDS(words []string) *WordListDS {
	var builder strings.Builder
	for i, word := range words {
		if i > 0 {
			builder.WriteByte('$') // Use '$' as a separator
		}
		builder.WriteString(word)
	}
	data := ([]byte)(builder.String())
	index := suffixarray.New(data)

	initalPrepend, initialAppend := getInitialGroups(data)

	return &WordListDS{data, index, initalPrepend, initialAppend}
}

func getInitialGroups(data []byte) (prependLetters []string, appendLetters []string) {
	wordsBytes := bytes.Split(data, []byte{'$'})
	prependSet := make(map[byte]bool)
	appendSet := make(map[byte]bool)
	for _, wordByte := range wordsBytes {
		for i := 0; i < len(wordByte); i++ {
			if i < len(wordByte)-1 {
				prependSet[wordByte[i]] = true
			}
			if i > 0 {
				appendSet[wordByte[i]] = true
			}
		}
	}
	// Convert to sorted slices
	for b := range prependSet {
		prependLetters = append(prependLetters, string(b))
	}
	sort.Strings(prependLetters)
	for b := range appendSet {
		appendLetters = append(appendLetters, string(b))
	}
	sort.Strings(appendLetters)
	return prependLetters, appendLetters
}

// GetGroups returns the prepend and append groups for a given query string S.
func (ds *WordListDS) GetGroups(S string) (prependLetters []string, appendLetters []string) {
	if len(S) == 0 {
		// Handle empty query
		return ds.initalPrepend, ds.initialAppend
	}

	positions := ds.index.Lookup([]byte(S), -1) // Get all positions of S
	prependSet := make(map[byte]bool)
	appendSet := make(map[byte]bool)

	for _, p := range positions {
		// Check prepend
		if p > 0 && ds.data[p-1] != '$' {
			prependSet[ds.data[p-1]] = true
		}
		// Check append
		end := p + len(S)
		if end < len(ds.data) && ds.data[end] != '$' {
			appendSet[ds.data[end]] = true
		}
	}

	// Convert sets to sorted slices
	var pre []string
	for b := range prependSet {
		pre = append(pre, string(b))
	}
	sort.Strings(pre)

	var app []string
	for b := range appendSet {
		app = append(app, string(b))
	}
	sort.Strings(app)

	return pre, app
}

// Contains checks if the given word exists in the word list.
func (ds *WordListDS) Contains(word string) bool {
	if len(word) == 0 {
		return false
	}
	// Look for the word followed by '$' or at the end of data
	pattern := word + "$"
	positions := ds.index.Lookup([]byte(pattern), -1)
	for _, p := range positions {
		// Ensure it's a full word by checking if it's at the start or preceded by '$'
		if p == 0 || ds.data[p-1] == '$' {
			return true
		}
	}
	// Check if the word is at the very end of the data
	if len(ds.data) >= len(word) && string(ds.data[len(ds.data)-len(word):]) == word {
		// Ensure it's a full word by checking if it's at the start or preceded by '$'
		if len(ds.data) == len(word) || ds.data[len(ds.data)-len(word)-1] == '$' {
			return true
		}
	}
	return false
}
