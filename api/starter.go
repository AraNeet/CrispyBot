package server

import (
	"CrispyBot/database"
	"CrispyBot/roller"
	"net/http"

	"github.com/labstack/echo/v4"
)

const PORT = ":8080"

func StartServer() {
	api := echo.New()

	// Root path
	api.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "CrispyBot API - A Discord bot with character generation")
	})

	// API paths
	api.GET("/api/health", healthCheck)

	// Character endpoints
	api.GET("/api/characters/:id", getCharacter)
	api.GET("/api/users/:id/character", getUserCharacter)
	api.POST("/api/users/:id/character", createCharacter)

	api.Logger.Fatal(api.Start(PORT))
}

// healthCheck returns the API health status
func healthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// getCharacter retrieves a character by ID
func getCharacter(ctx echo.Context) error {
	characterID := ctx.Param("id")
	if characterID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Character ID is required",
		})
	}

	// Connect to the database
	db := database.Connect()

	// Get the character
	character, err := database.GetCharacter(db, characterID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Character not found",
		})
	}

	return ctx.JSON(http.StatusOK, character)
}

// getUserCharacter retrieves a character by user ID (Discord ID)
func getUserCharacter(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}

	// Connect to the database
	db := database.Connect()

	// Get the character
	character, err := database.GetCharacterByOwner(db, userID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Character not found for this user",
		})
	}

	return ctx.JSON(http.StatusOK, character)
}

// createCharacter creates a new character for a user
func createCharacter(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}

	// Connect to the database
	db := database.Connect()

	// Check if user already has a character
	_, err := database.GetCharacterByOwner(db, userID)
	if err == nil {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"error": "User already has a character",
		})
	}

	// Generate a new character
	newCharacter := roller.GenerateCharacter(userID)

	// Save to the database
	savedChar, err := database.SaveCharacter(db, newCharacter)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create character: " + err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, savedChar)
}
