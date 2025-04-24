package models

import (
	"reflect"
	"testing"
)

// MockWordDictionary implements WordDictionaryI for testing.
type MockWordDictionary struct {
	words       map[string]bool
	prefixWords map[string][]string
	suffixWords map[string][]string
	wordList    []string
	forwardTrie *MockTrie
	reverseTrie *MockTrie
}

func (m *MockWordDictionary) ContainsWord(word string) bool {
	return m.words[word]
}
func (m *MockWordDictionary) FindWordsWithPrefix(prefix string) []string {
	return m.prefixWords[prefix]
}
func (m *MockWordDictionary) FindWordsWithSuffix(suffix string) []string {
	return m.suffixWords[suffix]
}

func (m *MockWordDictionary) GetForwardTrie() TrieI {
	return m.forwardTrie
}
func (m *MockWordDictionary) GetReverseTrie() TrieI {
	return m.reverseTrie
}

func (m *MockWordDictionary) GetWordList() []string {
	return m.wordList
}

// MockTrie for minimal interface usage in UpdateSets
type MockTrie struct {
	nextLetters  map[string][]string
	keysWithPref map[string][]string
}

func (t *MockTrie) GetNextLetters(prefix string) []string {
	return t.nextLetters[prefix]
}
func (t *MockTrie) KeysWithPrefix(prefix string) []string {
	return t.keysWithPref[prefix]
}

// --- Tests ---

func TestCheckValidWord(t *testing.T) {
	mockDict := &MockWordDictionary{
		words: map[string]bool{"cat": true, "dog": true},
	}
	tests := []struct {
		name  string
		state WordBuilderState
		want  bool
	}{
		{"Valid word", WordBuilderState{Answer: "cat"}, true},
		{"Invalid word", WordBuilderState{Answer: "bat"}, false},
		{"Empty answer", WordBuilderState{Answer: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckValidWord(tt.state, mockDict)
			if got != tt.want {
				t.Errorf("CheckValidWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddLetter(t *testing.T) {
	mockDict := &MockWordDictionary{
		words:       map[string]bool{"cat": true, "at": true},
		prefixWords: map[string][]string{"c": {"cat"}, "a": {"at"}},
		suffixWords: map[string][]string{"t": {"cat", "at"}},
		wordList:    []string{"cat", "at"},
		forwardTrie: &MockTrie{
			nextLetters:  map[string][]string{"ca": {"t"}},
			keysWithPref: map[string][]string{},
		},
		reverseTrie: &MockTrie{
			nextLetters:  map[string][]string{"ta": {"c"}},
			keysWithPref: map[string][]string{},
		},
	}
	state := WordBuilderState{
		Answer:    "at",
		PrefixSet: map[string]bool{"c": true},
		SuffixSet: map[string]bool{"t": true},
		Step:      0,
	}
	t.Run("Add valid prefix", func(t *testing.T) {
		newState, msg, err := AddLetter(state, mockDict, "c", "prefix")
		if err != nil || newState.Answer != "cat" || !newState.IsValidWord {
			t.Errorf("AddLetter() failed: %v, msg: %s, state: %+v", err, msg, newState)
		}
	})
	t.Run("Add invalid prefix", func(t *testing.T) {
		_, _, err := AddLetter(state, mockDict, "b", "prefix")
		if err == nil {
			t.Error("Expected error for invalid prefix letter")
		}
	})
	t.Run("Add invalid position", func(t *testing.T) {
		_, _, err := AddLetter(state, mockDict, "c", "middle")
		if err == nil {
			t.Error("Expected error for invalid position")
		}
	})
}

func TestRemoveLetter(t *testing.T) {
	forwardTrie := &MockTrie{
		nextLetters: map[string][]string{
			"a":  {"t"},
			"c":  {"a"},
			"ca": {"t"},
		},
	}
	reverseTrie := &MockTrie{
		nextLetters: map[string][]string{
			"ta": {"c"},
		},
		keysWithPref: map[string][]string{
			"ta": {"ta", "tac"},
		},
	}
	mockDict := &MockWordDictionary{
		words: map[string]bool{"cat": true, "at": true},
		prefixWords: map[string][]string{
			"a":   {"at"},
			"at":  {"at"},
			"c":   {"cat"},
			"ca":  {"cat"},
			"cat": {"cat"},
		},
		suffixWords: map[string][]string{
			"t": {"at", "cat"},
		},
		wordList:    []string{"cat", "at"},
		forwardTrie: forwardTrie,
		reverseTrie: reverseTrie,
	}
	state := WordBuilderState{
		Answer:    "at",
		PrefixSet: map[string]bool{"c": true},
		SuffixSet: map[string]bool{"t": true},
		Step:      0,
	}
	t.Run("Remove valid index", func(t *testing.T) {
		newState, _, err := RemoveLetter(state, mockDict, 0)
		if err != nil || newState.Answer != "t" {
			t.Errorf("RemoveLetter() failed: %v, state: %+v", err, newState)
		}
	})
	t.Run("Remove invalid index", func(t *testing.T) {
		_, _, err := RemoveLetter(state, mockDict, 5)
		if err == nil {
			t.Error("Expected error for out-of-bounds index")
		}
	})
}

func TestUpdateSets(t *testing.T) {
	mockDict := &MockWordDictionary{
		prefixWords: map[string][]string{"c": {"cat"}, "a": {"at"}},
		suffixWords: map[string][]string{"t": {"cat", "at"}},
		wordList:    []string{"cat", "at"},
		forwardTrie: &MockTrie{
			nextLetters:  map[string][]string{"ca": {"t"}},
			keysWithPref: map[string][]string{"ca": {"cat"}},
		},
		reverseTrie: &MockTrie{
			nextLetters:  map[string][]string{"ta": {"c"}},
			keysWithPref: map[string][]string{"ta": {"tac"}},
		},
	}
	state := WordBuilderState{
		Answer: "",
	}
	newState := UpdateSets(state, mockDict)
	if len(newState.PrefixSet) == 0 || len(newState.SuffixSet) == 0 {
		t.Error("UpdateSets() did not update prefix/suffix sets")
	}
}

func TestGetCurrentState(t *testing.T) {
	state := WordBuilderState{
		Answer:           "cat",
		PrefixSet:        map[string]bool{"a": true, "b": true},
		SuffixSet:        map[string]bool{"c": true},
		Step:             2,
		IsValidWord:      true,
		ValidCompletions: []string{"cat", "cats", "catch"},
		Suggestion:       "Try adding 's'",
	}
	got := GetCurrentState(state)
	if got["answer"] != "cat" || got["step"] != 2 || got["is_valid_word"] != true {
		t.Errorf("GetCurrentState() returned wrong values: %+v", got)
	}
	if !reflect.DeepEqual(got["prefix_set"], []string{"a", "b"}) && !reflect.DeepEqual(got["prefix_set"], []string{"b", "a"}) {
		t.Errorf("GetCurrentState() prefix_set mismatch: %+v", got["prefix_set"])
	}
}
