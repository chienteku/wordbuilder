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

// Store all EnhancedWordBuilder instances (session management)
var wordBuilders = make(map[string]*EnhancedWordBuilder)

// loadWordList loads the dictionary from a file
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
	// Load word list
	wordList, err := loadWordList("words.txt")
	if err != nil {
		log.Fatalf("Failed to load word list: %v", err)
	}
	fmt.Printf("Loaded %d words from words.txt\n", len(wordList))

	// Create the enhanced dictionary
	dictionary := NewWordDictionary(wordList)

	// Initialize Gin
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

	// Initialize WordBuilder
	r.POST("/api/wordbuilder/init", func(c *gin.Context) {
		sessionID := uuid.New().String()
		wordBuilders[sessionID] = NewEnhancedWordBuilder(dictionary)

		c.JSON(http.StatusOK, gin.H{
			"session_id": sessionID,
			"state":      wordBuilders[sessionID].GetCurrentState(),
			"success":    true,
		})
	})

	// Reset WordBuilder
	r.POST("/api/wordbuilder/reset", func(c *gin.Context) {
		var req struct {
			SessionID string `json:"session_id"`
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

		wb.Reset()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"state":   wb.GetCurrentState(),
			"message": "Word builder has been reset.",
		})
	})

	// Add letter
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

	// Query state
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

	// Start server
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
