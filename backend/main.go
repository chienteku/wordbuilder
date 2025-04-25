package main

import (
	"log"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"wordbuilder/controllers"
	"wordbuilder/models"
	"wordbuilder/services"
)

func main() {
	// Initialize database service
	dataDir := "data"
	uploadsDir := filepath.Join(dataDir, "uploads")
	dbPath := filepath.Join(dataDir, "wordbuilder.db")
	dbService, err := services.NewDatabaseService(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()

	// Initialize dictionary service
	dictService := services.NewDictionaryService()

	// Initialize word list service
	wordListService := services.NewWordListService(dbService, dictService, uploadsDir)

	// Load default word list if available
	wordLists, err := wordListService.GetAllWordLists()

	// Initialize dictionary - either from existing word list or default
	var dictionary *models.WordDictionary

	if err == nil && len(wordLists) > 0 {
		// Use the most recently updated word list
		dictionary, err = wordListService.LoadWordListIntoDictionary(wordLists[0].ID)
		if err != nil {
			log.Printf("Failed to load existing word list, falling back to default: %v", err)
		} else {
			log.Printf("Loaded dictionary from word list '%s' with %d words\n",
				wordLists[0].Name, len(dictionary.WordList))
		}
	}

	// Initialize services
	wordBuilderService := services.NewWordBuilderService(dictionary)

	// Initialize settings controller
	settingsController := controllers.NewSettingsController(dataDir, dbService)

	// Initialize controllers
	wordBuilderController := controllers.NewWordBuilderController(wordBuilderService)
	wordListController := controllers.NewWordListController(wordListService, wordBuilderService)
	dictionaryController := controllers.NewDictionaryController(settingsController)

	// Initialize Gin
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Register routes
	wordBuilderController.RegisterRoutes(r)
	wordListController.RegisterRoutes(r)
	dictionaryController.RegisterRoutes(r)
	settingsController.RegisterRoutes(r) // Register the settings routes

	// Start server
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
