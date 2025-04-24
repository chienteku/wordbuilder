package controllers

import (
	"net/http"
	"os"
	"wordbuilder/services"

	"github.com/gin-gonic/gin"
)

// SettingsController handles settings-related HTTP requests
type SettingsController struct {
	DBService *services.DatabaseService
}

// Settings represents the application settings
type Settings struct {
	PixabayAPIKey string `json:"pixabay_api_key"`
}

// NewSettingsController creates a new settings controller
func NewSettingsController(dataDir string, dbService *services.DatabaseService) *SettingsController {
	return &SettingsController{
		DBService: dbService,
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

	// Get the Pixabay API key from database
	pixabayKey, err := c.DBService.GetSetting("pixabay_api_key")
	if err == nil {
		settings.PixabayAPIKey = pixabayKey
	}

	return settings, nil
}

// saveSettings saves settings to the settings file
func (c *SettingsController) saveSettings(settings Settings) error {
	// Save Pixabay API key to database
	return c.DBService.SaveSetting("pixabay_api_key", settings.PixabayAPIKey)
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
