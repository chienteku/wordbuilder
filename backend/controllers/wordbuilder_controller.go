package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wordbuilder/services"
)

// WordBuilderController handles HTTP requests for wordbuilder game
type WordBuilderController struct {
	WordBuilderService *services.WordBuilderService
}

// NewWordBuilderController creates a new controller instance
func NewWordBuilderController(wbService *services.WordBuilderService) *WordBuilderController {
	return &WordBuilderController{
		WordBuilderService: wbService,
	}
}

// InitSession initializes a new WordBuilder session
func (c *WordBuilderController) InitSession(ctx *gin.Context) {
	sessionID := uuid.New().String()
	builder := c.WordBuilderService.CreateSession(sessionID)

	ctx.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"state":      builder.GetCurrentState(),
		"success":    true,
	})
}

// ResetSession resets a WordBuilder session
func (c *WordBuilderController) ResetSession(ctx *gin.Context) {
	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	builder, exists := c.WordBuilderService.ResetSession(req.SessionID)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"state":   builder.GetCurrentState(),
		"message": "Word builder has been reset.",
	})
}

// AddLetter adds a letter to the word
func (c *WordBuilderController) AddLetter(ctx *gin.Context) {
	var req struct {
		SessionID string `json:"session_id"`
		Letter    string `json:"letter"`
		Position  string `json:"position"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	builder, exists := c.WordBuilderService.GetSession(req.SessionID)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if len(req.Letter) != 1 || !strings.Contains("abcdefghijklmnopqrstuvwxyz", strings.ToLower(req.Letter)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Letter must be a single lowercase letter"})
		return
	}
	if req.Position != "prefix" && req.Position != "suffix" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Position must be 'prefix' or 'suffix'"})
		return
	}

	success, message := builder.AddLetter(strings.ToLower(req.Letter), req.Position)
	if !success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"state":   builder.GetCurrentState(),
		"message": message,
	})
}

// RemoveLetter removes a letter from the word
func (c *WordBuilderController) RemoveLetter(ctx *gin.Context) {
	var req struct {
		SessionID string `json:"session_id"`
		Index     int    `json:"index"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	builder, exists := c.WordBuilderService.GetSession(req.SessionID)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	success, message := builder.RemoveLetter(req.Index)
	if !success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"state":   builder.GetCurrentState(),
		"message": message,
	})
}

// GetState returns the current state of a session
func (c *WordBuilderController) GetState(ctx *gin.Context) {
	sessionID := ctx.Query("session_id")
	builder, exists := c.WordBuilderService.GetSession(sessionID)
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"state": builder.GetCurrentState(),
	})
}

// RegisterRoutes registers all controller routes
func (c *WordBuilderController) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/wordbuilder")
	{
		api.POST("/init", c.InitSession)
		api.POST("/reset", c.ResetSession)
		api.POST("/add", c.AddLetter)
		api.POST("/remove", c.RemoveLetter)
		api.GET("/state", c.GetState)
	}
}
