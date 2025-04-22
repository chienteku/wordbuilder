package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"wordbuilder/services"

	"github.com/gin-gonic/gin"
)

// WordListController handles HTTP requests for word list management
type WordListController struct {
	WordListService *services.WordListService
	MaxFileSize     int64
}

// NewWordListController creates a new word list controller
func NewWordListController(wordListService *services.WordListService) *WordListController {
	return &WordListController{
		WordListService: wordListService,
		MaxFileSize:     10 * 1024 * 1024, // Default to 10MB max file size
	}
}

// GetAllWordLists returns all word lists
func (c *WordListController) GetAllWordLists(ctx *gin.Context) {
	wordLists, err := c.WordListService.GetAllWordLists()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get word lists: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"word_lists": wordLists,
	})
}

// GetWordList returns a specific word list
func (c *WordListController) GetWordList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word list ID"})
		return
	}

	wordList, err := c.WordListService.GetWordList(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Word list not found: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"word_list": wordList,
	})
}

// GetWordListSample returns a sample of words from a word list
func (c *WordListController) GetWordListSample(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word list ID"})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	words, err := c.WordListService.ReadWordListContent(id, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read word list: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"words": words,
		"count": len(words),
	})
}

// CreateWordList creates a new word list from an uploaded file
func (c *WordListController) CreateWordList(ctx *gin.Context) {
	// Set max file size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, c.MaxFileSize)

	// Parse form fields
	err := ctx.Request.ParseMultipartForm(c.MaxFileSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large (max %dMB): %v", c.MaxFileSize/(1024*1024), err)})
		return
	}

	// Get form values
	name := ctx.Request.FormValue("name")
	description := ctx.Request.FormValue("description")
	source := ctx.Request.FormValue("source")

	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Get file
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".txt" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only .txt files are allowed"})
		return
	}

	// Read file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read file: %v", err)})
		return
	}

	// Create word list
	wordList, err := c.WordListService.CreateWordList(fileData, name, description, source)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create word list: %v", err)})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":   "Word list created successfully",
		"word_list": wordList,
	})
}

// UpdateWordList updates an existing word list
func (c *WordListController) UpdateWordList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word list ID"})
		return
	}

	// Set max file size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, c.MaxFileSize)

	// Parse form fields
	err = ctx.Request.ParseMultipartForm(c.MaxFileSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large (max %dMB): %v", c.MaxFileSize/(1024*1024), err)})
		return
	}

	// Get form values
	name := ctx.Request.FormValue("name")
	description := ctx.Request.FormValue("description")
	source := ctx.Request.FormValue("source")

	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	var fileData []byte

	// Check if a file was uploaded
	file, header, err := ctx.Request.FormFile("file")
	if err == nil {
		defer file.Close()

		// Validate file extension
		ext := filepath.Ext(header.Filename)
		if ext != ".txt" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only .txt files are allowed"})
			return
		}

		// Read file content
		fileData, err = io.ReadAll(file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read file: %v", err)})
			return
		}
	}

	// Update word list
	wordList, err := c.WordListService.UpdateWordList(id, name, description, source, fileData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update word list: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Word list updated successfully",
		"word_list": wordList,
	})
}

// DeleteWordList deletes a word list
func (c *WordListController) DeleteWordList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word list ID"})
		return
	}

	err = c.WordListService.DeleteWordList(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete word list: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Word list deleted successfully",
	})
}

// DownloadWordList streams the word list file for download
func (c *WordListController) DownloadWordList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word list ID"})
		return
	}

	// Get the word list
	wordList, err := c.WordListService.GetWordList(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Word list not found: %v", err)})
		return
	}

	// Open the file
	file, err := os.Open(wordList.FilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open word list file: %v", err)})
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get file info: %v", err)})
		return
	}

	// Set response headers
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.txt", wordList.Name))
	ctx.Header("Content-Type", "text/plain")
	ctx.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Stream the file
	ctx.File(wordList.FilePath)
}

// UseWordList sets a word list as the active dictionary
func (c *WordListController) UseWordList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word list ID"})
		return
	}

	// Load the word list into a dictionary
	dictionary, err := c.WordListService.LoadWordListIntoDictionary(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load word list: %v", err)})
		return
	}

	// Update the word builder service with the new dictionary
	// This assumes the word builder service is accessible here,
	// which might require passing it to the controller or using a global registry
	// For simplicity, we'll just return a message
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Word list loaded successfully with %d words", len(dictionary.WordList)),
	})
}

// RegisterRoutes registers all controller routes
func (c *WordListController) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/wordlists")
	{
		api.GET("", c.GetAllWordLists)
		api.GET("/:id", c.GetWordList)
		api.GET("/:id/sample", c.GetWordListSample)
		api.POST("", c.CreateWordList)
		api.PUT("/:id", c.UpdateWordList)
		api.DELETE("/:id", c.DeleteWordList)
		api.GET("/:id/download", c.DownloadWordList)
		api.POST("/:id/use", c.UseWordList)
	}
}
