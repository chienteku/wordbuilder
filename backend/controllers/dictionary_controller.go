package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

// DictionaryApiResponse represents the response from the Dictionary API
type DictionaryApiResponse []struct {
	Word      string `json:"word"`
	Phonetics []struct {
		Text  string `json:"text,omitempty"`
		Audio string `json:"audio,omitempty"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example,omitempty"`
		} `json:"definitions"`
	} `json:"meanings"`
}

// WordDetails represents word details we want to return to the frontend
type WordDetails struct {
	Pronunciation string `json:"pronunciation"`
	Audio         string `json:"audio"`
	Meaning       string `json:"meaning"`
	Example       string `json:"example"`
	ImageUrl      string `json:"imageUrl,omitempty"`
}

// DictionaryController handles dictionary-related requests
type DictionaryController struct {
	SettingsController *SettingsController
}

// NewDictionaryController creates a new dictionary controller
func NewDictionaryController(settingsController *SettingsController) *DictionaryController {
	return &DictionaryController{
		SettingsController: settingsController,
	}
}

// GetWordImage fetches an image for a word from Pixabay
func (c *DictionaryController) GetWordImage(ctx *gin.Context) {
	word := ctx.Param("word")
	if word == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Word parameter is required"})
		return
	}

	// Get API key from settings
	apiKey := c.SettingsController.GetPixabayAPIKey()
	if apiKey == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Pixabay API key not configured"})
		return
	}

	// Build Pixabay API URL
	url := fmt.Sprintf("https://pixabay.com/api/?key=%s&q=%s&image_type=photo", apiKey, word)

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

// GetWordDetails fetches word details from the Dictionary API
func (c *DictionaryController) GetWordDetails(ctx *gin.Context) {
	word := ctx.Param("word")
	if word == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Word parameter is required"})
		return
	}

	// Build Dictionary API URL
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)

	// Make the request to Dictionary API
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch word details"})
		return
	}
	defer resp.Body.Close()

	// Check if the word was found
	if resp.StatusCode == http.StatusNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Word not found in dictionary"})
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Parse the response
	var dictResp DictionaryApiResponse
	if err := json.Unmarshal(body, &dictResp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse dictionary response"})
		return
	}

	// Extract the details we need
	details := WordDetails{}

	// Check if we got valid data
	if len(dictResp) > 0 {
		// Get pronunciation
		if len(dictResp[0].Phonetics) > 0 {
			details.Pronunciation = dictResp[0].Phonetics[0].Text
			details.Audio = dictResp[0].Phonetics[0].Audio
		}

		// Get meaning
		if len(dictResp[0].Meanings) > 0 && len(dictResp[0].Meanings[0].Definitions) > 0 {
			details.Meaning = dictResp[0].Meanings[0].Definitions[0].Definition
			details.Example = dictResp[0].Meanings[0].Definitions[0].Example
		}
	}

	// Return the word details
	ctx.JSON(http.StatusOK, details)
}

// GetCompleteWordDetails fetches both definition and image for a word
func (c *DictionaryController) GetCompleteWordDetails(ctx *gin.Context) {
	word := ctx.Param("word")
	if word == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Word parameter is required"})
		return
	}

	// Get word details from dictionary API
	details := WordDetails{}

	// Build Dictionary API URL
	dictUrl := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)

	// Make the request to Dictionary API
	dictResp, err := http.Get(dictUrl)
	if err == nil && dictResp.StatusCode == http.StatusOK {
		defer dictResp.Body.Close()

		// Read and parse response
		body, err := io.ReadAll(dictResp.Body)
		if err == nil {
			var apiResp DictionaryApiResponse
			if json.Unmarshal(body, &apiResp) == nil && len(apiResp) > 0 {
				// Get pronunciation
				if len(apiResp[0].Phonetics) > 0 {
					details.Pronunciation = apiResp[0].Phonetics[0].Text
					details.Audio = apiResp[0].Phonetics[0].Audio
				}

				// Get meaning
				if len(apiResp[0].Meanings) > 0 && len(apiResp[0].Meanings[0].Definitions) > 0 {
					details.Meaning = apiResp[0].Meanings[0].Definitions[0].Definition
					details.Example = apiResp[0].Meanings[0].Definitions[0].Example
				}
			}
		}
	}

	// Get image from Pixabay
	apiKey := c.SettingsController.GetPixabayAPIKey()
	if apiKey != "" {
		// Build Pixabay API URL
		imgUrl := fmt.Sprintf("https://pixabay.com/api/?key=%s&q=%s&image_type=photo", apiKey, word)

		// Make the request to Pixabay
		imgResp, err := http.Get(imgUrl)
		if err == nil {
			defer imgResp.Body.Close()

			// Read and parse response
			body, err := io.ReadAll(imgResp.Body)
			if err == nil {
				var pixabayResp PixabayResponse
				if json.Unmarshal(body, &pixabayResp) == nil && len(pixabayResp.Hits) > 0 {
					details.ImageUrl = pixabayResp.Hits[0].WebformatURL
				}
			}
		}
	}

	// Return the complete word details
	ctx.JSON(http.StatusOK, details)
}

// RegisterRoutes registers all controller routes
func (c *DictionaryController) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/dictionary")
	{
		api.GET("/image/:word", c.GetWordImage)
		api.GET("/details/:word", c.GetWordDetails)
		api.GET("/complete/:word", c.GetCompleteWordDetails)
	}
}
