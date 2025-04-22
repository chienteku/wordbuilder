package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"wordbuilder/controllers"
	"wordbuilder/services"
)

func main() {
	// Initialize services
	dictService := services.NewDictionaryService()

	// Load word list
	wordList, err := dictService.LoadWordList("words.txt")
	if err != nil {
		log.Fatalf("Failed to load word list: %v", err)
	}
	fmt.Printf("Loaded %d words from words.txt\n", len(wordList))

	// Create the dictionary
	dictionary := dictService.CreateDictionary(wordList)

	// Initialize services
	wordBuilderService := services.NewWordBuilderService(dictionary)

	// Initialize controllers
	wordBuilderController := controllers.NewWordBuilderController(wordBuilderService)

	// Initialize Gin
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Register routes
	wordBuilderController.RegisterRoutes(r)

	// Start server
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
