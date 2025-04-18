package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TrieNode 表示 Trie 的節點
type TrieNode struct {
	Children map[rune]*TrieNode
	IsWord   bool
}

// Trie 表示單字庫的 Trie 結構
type Trie struct {
	Root *TrieNode
}

// NewTrie 創建新的 Trie
func NewTrie() *Trie {
	return &Trie{Root: &TrieNode{Children: make(map[rune]*TrieNode)}}
}

// Insert 插入單字到 Trie
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

// Contains 檢查單字是否存在於 Trie
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

// KeysWithPrefix 返回以指定前綴開頭的所有單字
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

// collectKeys 遞迴收集以當前節點為根的單字
func (t *Trie) collectKeys(node *TrieNode, prefix string, results *[]string) {
	if node.IsWord {
		*results = append(*results, prefix)
	}
	for ch, child := range node.Children {
		t.collectKeys(child, prefix+string(ch), results)
	}
}

// WordBuilder 表示單字構建器
type WordBuilder struct {
	Answer      string
	PrefixSet   map[string]bool
	SuffixSet   map[string]bool
	Step        int
	Trie        *Trie
	IsValidWord bool // Added isValidWord field
}

// NewWordBuilder 創建新的 WordBuilder
func NewWordBuilder(trie *Trie) *WordBuilder {
	wb := &WordBuilder{
		Answer:      "",
		PrefixSet:   make(map[string]bool),
		SuffixSet:   make(map[string]bool),
		Step:        0,
		Trie:        trie,
		IsValidWord: false, // Initialize as false
	}

	// 初始時提供所有26個字母
	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		wb.PrefixSet[string(letter)] = true
		wb.SuffixSet[string(letter)] = true
	}

	return wb
}

// CheckValidWord 檢查當前答案是否為有效單字
func (wb *WordBuilder) CheckValidWord() bool {
	isValid := len(wb.Answer) > 0 && wb.Trie.Contains(wb.Answer)
	wb.IsValidWord = isValid // Update the IsValidWord field
	return isValid
}

// AddLetter 添加字母到答案
func (wb *WordBuilder) AddLetter(letter, position string) (bool, string) {
	if position == "prefix" {
		if !wb.PrefixSet[letter] {
			return false, fmt.Sprintf("Invalid letter '%s' for prefix set.", letter)
		}
		wb.Answer = letter + wb.Answer
	} else if position == "suffix" {
		if !wb.SuffixSet[letter] {
			return false, fmt.Sprintf("Invalid letter '%s' for suffix set.", letter)
		}
		wb.Answer = wb.Answer + letter
	} else {
		return false, "Invalid position. Use 'prefix' or 'suffix'."
	}

	// Update IsValidWord field with current word validity
	wb.CheckValidWord()

	message := fmt.Sprintf("Step %d: Added '%s' as %s -> Answer: %s", wb.Step, letter, position, wb.Answer)
	if wb.IsValidWord {
		message += fmt.Sprintf("\n*** '%s' is a valid English word! ***", wb.Answer)
	}

	wb.UpdateSets()
	wb.Step++
	return true, message
}

// RemoveLetter removes a letter at the specified index
func (wb *WordBuilder) RemoveLetter(index int) (bool, string) {
	if index < 0 || index >= len(wb.Answer) {
		return false, fmt.Sprintf("Invalid index %d for answer '%s'.", index, wb.Answer)
	}

	// Allow removing any letter, including the last one
	letter := string(wb.Answer[index])
	newAnswer := wb.Answer[:index] + wb.Answer[index+1:]
	wb.Answer = newAnswer

	// Update IsValidWord field with current word validity
	wb.CheckValidWord()

	message := fmt.Sprintf("Step %d: Removed '%s' at index %d -> Answer: %s", wb.Step, letter, index, wb.Answer)
	if wb.IsValidWord {
		message += fmt.Sprintf("\n*** '%s' is a valid English word! ***", wb.Answer)
	}

	wb.UpdateSets()
	wb.Step++
	return true, message
}

// UpdateSets 動態生成前綴和後綴字母集合
func (wb *WordBuilder) UpdateSets() {
	wb.PrefixSet = make(map[string]bool)
	wb.SuffixSet = make(map[string]bool)

	// 如果沒有字母，則提供所有26個字母
	if len(wb.Answer) == 0 {
		for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
			wb.PrefixSet[string(letter)] = true
			wb.SuffixSet[string(letter)] = true
		}
		return
	}

	// 生成後綴集合
	words := wb.Trie.KeysWithPrefix(wb.Answer)
	for _, word := range words {
		if len(word) > len(wb.Answer) {
			nextLetter := string(word[len(wb.Answer)])
			wb.SuffixSet[nextLetter] = true
		}
	}

	// 生成前綴集合
	for _, letter := range "abcdefghijklmnopqrstuvwxyz" {
		testPrefix := string(letter) + wb.Answer
		if len(wb.Trie.KeysWithPrefix(testPrefix)) > 0 {
			wb.PrefixSet[string(letter)] = true
		}
	}
}

// GetCurrentState 返回當前狀態
func (wb *WordBuilder) GetCurrentState() map[string]interface{} {
	prefixSet := make([]string, 0, len(wb.PrefixSet))
	for letter := range wb.PrefixSet {
		prefixSet = append(prefixSet, letter)
	}
	suffixSet := make([]string, 0, len(wb.SuffixSet))
	for letter := range wb.SuffixSet {
		suffixSet = append(suffixSet, letter)
	}
	return map[string]interface{}{
		"answer":        wb.Answer,
		"prefix_set":    prefixSet,
		"suffix_set":    suffixSet,
		"step":          wb.Step,
		"is_valid_word": wb.IsValidWord, // Include is_valid_word in the state
	}
}

// 儲存所有 WordBuilder 實例（模擬 session）
var wordBuilders = make(map[string]*WordBuilder)

// loadWordList 從檔案載入單字庫
func loadWordList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if word != "" {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

func main() {
	// 載入單字庫
	wordList, err := loadWordList("words.txt")
	if err != nil {
		log.Fatalf("Failed to load word list: %v", err)
	}
	fmt.Printf("Loaded %d words from words.txt\n", len(wordList))

	// 構建 Trie
	trie := NewTrie()
	for _, word := range wordList {
		trie.Insert(word)
	}

	// 初始化 Gin
	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Replace with your frontend's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// 初始化 WordBuilder
	r.POST("/api/wordbuilder/init", func(c *gin.Context) {
		sessionID := uuid.New().String()
		wordBuilders[sessionID] = NewWordBuilder(trie)

		c.JSON(http.StatusOK, gin.H{
			"session_id": sessionID,
			"state":      wordBuilders[sessionID].GetCurrentState(),
			"success":    true,
		})
	})

	// 添加字母
	r.POST("/api/wordbuilder/add", func(c *gin.Context) {
		var req struct {
			SessionID string `json:"session_id"`
			Letter    string `json:"letter"`
			Position  string `json:"position"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		wb, exists := wordBuilders[req.SessionID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}

		if len(req.Letter) != 1 || !strings.Contains("abcdefghijklmnopqrstuvwxyz", strings.ToLower(req.Letter)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Letter must be a single lowercase letter"})
			return
		}
		if req.Position != "prefix" && req.Position != "suffix" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Position must be 'prefix' or 'suffix'"})
			return
		}

		success, message := wb.AddLetter(strings.ToLower(req.Letter), req.Position)
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"state":   wb.GetCurrentState(),
			"message": message,
		})
	})

	// Remove letter
	r.POST("/api/wordbuilder/remove", func(c *gin.Context) {
		var req struct {
			SessionID string `json:"session_id"`
			Index     int    `json:"index"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		wb, exists := wordBuilders[req.SessionID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}

		success, message := wb.RemoveLetter(req.Index)
		if !success {
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"state":   wb.GetCurrentState(),
			"message": message,
		})
	})

	// 查詢狀態
	r.GET("/api/wordbuilder/state", func(c *gin.Context) {
		sessionID := c.Query("session_id")
		wb, exists := wordBuilders[sessionID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"state": wb.GetCurrentState(),
		})
	})

	// 啟動伺服器
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
