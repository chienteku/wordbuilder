package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SettingsController handles settings-related HTTP requests
type SettingsController struct {
	SettingsFilePath string
}

// Settings represents the application settings
type Settings struct {
	PixabayAPIKey string `json:"pixabay_api_key"`
}

// NewSettingsController creates a new settings controller
func NewSettingsController(dataDir string) *SettingsController {
	// Ensure the settings file path exists
	settingsDir := filepath.Join(dataDir, "settings")
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		// Log the error but continue - we'll handle file creation when saving
		// In a production app, you might want to handle this more gracefully
		println("Warning: Failed to create settings directory:", err.Error())
	}

	return &SettingsController{
		SettingsFilePath: filepath.Join(settingsDir, "settings.json"),
	}
}

// GetSettings retrieves the current settings
func (c *SettingsController) GetSettings(ctx *gin.Context) {
	settings, err := c.loadSettings()
	if err != nil {
		// If settings file doesn't exist, return empty settings
		settings = Settings{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"settings": settings,
	})
}

// UpdateSettings updates the application settings
func (c *SettingsController) UpdateSettings(ctx *gin.Context) {
	var settings Settings
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings data"})
		return
	}

	// Save settings to file
	if err := c.saveSettings(settings); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Settings updated successfully",
		"settings": settings,
	})
}

// loadSettings loads settings from the settings file
func (c *SettingsController) loadSettings() (Settings, error) {
	var settings Settings

	// Check if the file exists
	if _, err := os.Stat(c.SettingsFilePath); os.IsNotExist(err) {
		return settings, err
	}

	// Read the file
	data, err := os.ReadFile(c.SettingsFilePath)
	if err != nil {
		return settings, err
	}

	// Parse the JSON
	if err := json.Unmarshal(data, &settings); err != nil {
		return settings, err
	}

	return settings, nil
}

// saveSettings saves settings to the settings file
func (c *SettingsController) saveSettings(settings Settings) error {
	// Convert settings to JSON
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(c.SettingsFilePath, data, 0644)
}

// GetPixabayAPIKey returns the Pixabay API key
func (c *SettingsController) GetPixabayAPIKey() string {
	// Try to load from settings file first
	settings, err := c.loadSettings()
	if err == nil && settings.PixabayAPIKey != "" {
		return settings.PixabayAPIKey
	}

	// Fall back to environment variable
	return os.Getenv("PIXABAY_API_KEY")
}

// RegisterRoutes registers all controller routes
func (c *SettingsController) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/settings")
	{
		api.GET("", c.GetSettings)
		api.POST("", c.UpdateSettings)
	}
}
