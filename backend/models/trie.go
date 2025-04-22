package models

// TrieNode represents a node in the Trie
type TrieNode struct {
	Children map[rune]*TrieNode
	IsWord   bool
}

// Trie represents a dictionary structure for efficient prefix searches
type Trie struct {
	Root *TrieNode
}

// NewTrie creates a new Trie
func NewTrie() *Trie {
	return &Trie{Root: &TrieNode{Children: make(map[rune]*TrieNode)}}
}

// Insert adds a word to the Trie
func (t *Trie) Insert(word string) {
	node := t.Root
	for _, ch := range word {
		if _, exists := node.Children[ch]; !exists {
			node.Children[ch] = &TrieNode{Children: make(map[rune]*TrieNode)}
		}
		node = node.Children[ch]
	}
	node.IsWord = true
}

// Contains checks if a word exists in the Trie
func (t *Trie) Contains(word string) bool {
	node := t.Root
	for _, ch := range word {
		if _, exists := node.Children[ch]; !exists {
			return false
		}
		node = node.Children[ch]
	}
	return node.IsWord
}

// KeysWithPrefix returns all words that start with the given prefix
func (t *Trie) KeysWithPrefix(prefix string) []string {
	var results []string
	node := t.Root
	for _, ch := range prefix {
		if _, exists := node.Children[ch]; !exists {
			return results
		}
		node = node.Children[ch]
	}
	t.collectKeys(node, prefix, &results)
	return results
}

// collectKeys recursively collects all words from the current node
func (t *Trie) collectKeys(node *TrieNode, prefix string, results *[]string) {
	if node.IsWord {
		*results = append(*results, prefix)
	}
	for ch, child := range node.Children {
		t.collectKeys(child, prefix+string(ch), results)
	}
}

// GetNextLetters returns the possible next letters after a given prefix
func (t *Trie) GetNextLetters(prefix string) []string {
	node := t.Root
	for _, ch := range prefix {
		if child, exists := node.Children[ch]; exists {
			node = child
		} else {
			return []string{}
		}
	}
	letters := make([]string, 0, len(node.Children))
	for ch := range node.Children {
		letters = append(letters, string(ch))
	}
	return letters
}
