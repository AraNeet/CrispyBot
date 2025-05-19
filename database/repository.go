package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection names
const (
	usersCollection      = "users"
	charactersCollection = "characters"
)

// CreateUser creates a new user in the database
func CreateUser(db *DB, userID string) (User, error) {
	if db == nil {
		return User{}, fmt.Errorf("database connection is nil")
	}

	if userID == "" {
		return User{}, fmt.Errorf("discord ID isn't passed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	collection := db.GetCollection(usersCollection)
	filter := bson.M{"discordID": userID}

	var existingUser User
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		// User already exists
		return existingUser, nil
	} else if err != mongo.ErrNoDocuments {
		// Unexpected error
		return User{}, fmt.Errorf("error checking for existing user: %w", err)
	}

	// Create new user
	newUser := User{
		DiscordID: userID,
	}

	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Set the ID from the inserted document
	newUser.ID = result.InsertedID.(primitive.ObjectID)

	return newUser, nil
}

// GetUserByID retrieves a user by Discord ID
func GetUserByID(db *DB, discordID string) (User, error) {
	if db == nil {
		return User{}, fmt.Errorf("database connection is nil")
	}

	if discordID == "" {
		return User{}, fmt.Errorf("discord ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(usersCollection)
	filter := bson.M{"discordID": discordID}

	var user User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// SaveCharacter saves a character to the database
func SaveCharacter(db *DB, character Character, discordID string) (Character, error) {
	if db == nil {
		return Character{}, fmt.Errorf("database connection is nil")
	}

	if discordID == "" {
		return Character{}, fmt.Errorf("discord ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Associate character with owner
	character.Owner = discordID

	// Insert character
	collection := db.GetCollection(charactersCollection)
	result, err := collection.InsertOne(ctx, character)
	if err != nil {
		return Character{}, fmt.Errorf("failed to save character: %w", err)
	}

	// Get the ID of the inserted document
	character.ID = result.InsertedID.(primitive.ObjectID)

	// Also update the user document to reference this character
	userCollection := db.GetCollection(usersCollection)
	filter := bson.M{"discordID": discordID}
	update := bson.M{"$set": bson.M{"character": character}}
	opts := options.Update().SetUpsert(true)

	_, err = userCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return Character{}, fmt.Errorf("failed to update user with character: %w", err)
	}

	return character, nil
}

// GetCharacterByOwner retrieves a character by owner ID
func GetCharacterByOwner(db *DB, ownerID string) (Character, error) {
	if db == nil {
		return Character{}, fmt.Errorf("database connection is nil")
	}

	if ownerID == "" {
		return Character{}, fmt.Errorf("owner ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(charactersCollection)
	filter := bson.M{"owner": ownerID}

	var character Character
	err := collection.FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Character{}, fmt.Errorf("no character found for user: %s", ownerID)
		}
		return Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	return character, nil
}

// GetCharacter retrieves a character by ID
func GetCharacter(db *DB, characterID string) (Character, error) {
	if db == nil {
		return Character{}, fmt.Errorf("database connection is nil")
	}

	if characterID == "" {
		return Character{}, fmt.Errorf("character ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return Character{}, fmt.Errorf("invalid character ID format: %w", err)
	}

	collection := db.GetCollection(charactersCollection)
	filter := bson.M{"_id": objectID}

	var character Character
	err = collection.FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Character{}, fmt.Errorf("character not found")
		}
		return Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	return character, nil
}
