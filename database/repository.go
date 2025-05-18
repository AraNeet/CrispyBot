package database

import (
	"fmt"

	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func CreateUser(db *surrealdb.DB, userID string) (User, error) {
	if userID == "" {
		return User{}, fmt.Errorf("discord ID isn't passed.")
	}

	user, err := surrealdb.Create[User](db, models.Table("users"), map[interface{}]interface{}{
		"DiscordID": userID,
	})
	if err != nil {
		return User{}, fmt.Errorf("Failed to create user")
	}

	return *user, nil
}

// SaveCharacter saves a character to the database
func SaveCharacter(db *surrealdb.DB, character Character, discordID string) (Character, error) {
	if character.ID == "" {
		return Character{}, fmt.Errorf("character owner ID is required")
	}

	// Create the character record
	savedChar, err := surrealdb.Create[Character](db, models.Table("characters"), map[interface{}]interface{}{
		"owner":      discordID,
		"stats":      character.Stats,
		"attributes": character.Attributes,
	})
	if err != nil {
		return Character{}, fmt.Errorf("failed to save character: %w", err)
	}
	return *savedChar, nil
}

// GetCharacterByOwner retrieves a character by owner ID
func GetCharacterByOwner(db *surrealdb.DB, ownerID string) (Character, error) {
	if ownerID == "" {
		return Character{}, fmt.Errorf("owner ID is required")
	}

	// Query the database for characters with the given owner ID
	results, err := surrealdb.Select[Character](db, ownerID)
	if err != nil {
		return Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	character := Character{
		ID:         results.ID,
		Stats:      results.Stats,
		Attributes: results.Attributes,
	}

	return character, nil
}

// GetUserByID retrieves a user by Discord ID
func GetUserByID(db *surrealdb.DB, discordID string) (User, error) {
	if discordID == "" {
		return User{}, fmt.Errorf("discord ID is required")
	}

	// Query the database for the user with the given Discord ID
	user, err := surrealdb.Select[User](db, discordID)
	if err != nil {
		return User{}, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return User{}, fmt.Errorf("user not found")
	}

	return *user, nil
}
