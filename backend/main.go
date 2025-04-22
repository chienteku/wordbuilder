package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"wordbuilder/controllers"
	"wordbuilder/models"
	"wordbuilder/services"
)

func main() {
	// Create necessary directories
	dataDir := "data"
	uploadsDir := filepath.Join(dataDir, "uploads")

	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("Failed to create data directories: %v", err)
	}

	// Initialize database service
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

	// Fall back to default dictionary if needed
	if dictionary == nil {
		defaultWordListPath := "words.txt"
		wordList, err := dictService.LoadWordList(defaultWordListPath)
		if err != nil {
			log.Fatalf("Failed to load default word list: %v", err)
		}
		dictionary = dictService.CreateDictionary(wordList)
		log.Printf("Loaded %d words from default dictionary\n", len(dictionary.WordList))

		// Save the default dictionary as a word list if we have none
		if len(wordLists) == 0 {
			// Copy the default word list to the uploads directory
			defaultData, err := os.ReadFile(defaultWordListPath)
			if err == nil {
				wordListService.CreateWordList(
					defaultData,
					"Default Word List",
					"System default word list loaded at startup",
					"system",
				)
			}
		}
	}

	// Initialize services
	wordBuilderService := services.NewWordBuilderService(dictionary)

	// Initialize controllers
	wordBuilderController := controllers.NewWordBuilderController(wordBuilderService)
	wordListController := controllers.NewWordListController(wordListService, wordBuilderService) // Pass wordBuilderService to the controller

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

	// Start server
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
