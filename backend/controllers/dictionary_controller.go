package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// PixabayResponse represents the Pixabay API response structure
type PixabayResponse struct {
	Total     int          `json:"total"`
	TotalHits int          `json:"totalHits"`
	Hits      []PixabayHit `json:"hits"`
}

// PixabayHit represents a single image result from Pixabay
type PixabayHit struct {
	WebformatURL string `json:"webformatURL"`
	// Add other fields if needed
}

// DictionaryController handles dictionary-related requests
type DictionaryController struct {
	PixabayAPIKey string
}

// NewDictionaryController creates a new dictionary controller
func NewDictionaryController() *DictionaryController {
	// Get API key from environment variable
	apiKey := os.Getenv("PIXABAY_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: PIXABAY_API_KEY environment variable not set")
	}

	return &DictionaryController{
		PixabayAPIKey: apiKey,
	}
}

// GetWordImage fetches an image for a word from Pixabay
func (c *DictionaryController) GetWordImage(ctx *gin.Context) {
	word := ctx.Param("word")
	if word == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Word parameter is required"})
		return
	}

	// Build Pixabay API URL
	// fmt.Sprintf("https://pixabay.com/api/?key=API_KEY&q=answer&image_type=photo")
	url := fmt.Sprintf("https://pixabay.com/api/?key=%s&q=%s&image_type=photo", c.PixabayAPIKey, word)

	// Make the request to Pixabay
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Parse the response
	var pixabayResp PixabayResponse
	if err := json.Unmarshal(body, &pixabayResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Check if we got any hits
	if len(pixabayResp.Hits) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No images found for this word"})
		return
	}

	// Return the first image URL
	ctx.JSON(http.StatusOK, gin.H{
		"imageUrl": pixabayResp.Hits[0].WebformatURL,
	})
}

// RegisterRoutes registers all controller routes
func (c *DictionaryController) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/dictionary")
	{
		api.GET("/image/:word", c.GetWordImage)
	}
}
