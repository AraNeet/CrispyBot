package database

import (
	"CrispyBot/database/models"
	"CrispyBot/roller"
	"CrispyBot/variables"
	"context"
	"fmt"
	"math/rand"
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
func CreateUser(db *DB, userID string) (models.User, error) {
	if db == nil {
		return models.User{}, fmt.Errorf("database connection is nil")
	}

	if userID == "" {
		return models.User{}, fmt.Errorf("discord ID isn't passed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	collection := db.GetCollection(usersCollection)
	filter := bson.M{"discordID": userID}

	var existingUser models.User
	err := collection.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		// User already exists
		return existingUser, nil
	} else if err != mongo.ErrNoDocuments {
		// Unexpected error
		return models.User{}, fmt.Errorf("error checking for existing user: %w", err)
	}

	// Create new user with initial reroll counts
	newUser := models.User{
		DiscordID:       userID,
		Wallet:          0,
		FullRerolls:     2,
		StatRerolls:     1,
		LastRerollReset: time.Now(),
	}

	result, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Set the ID from the inserted document
	newUser.ID = result.InsertedID.(primitive.ObjectID)

	return newUser, nil
}

// GetUserByID retrieves a user by Discord ID
func GetUserByID(db *DB, discordID string) (models.User, error) {
	if db == nil {
		return models.User{}, fmt.Errorf("database connection is nil")
	}

	if discordID == "" {
		return models.User{}, fmt.Errorf("discord ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(usersCollection)
	filter := bson.M{"discordID": discordID}

	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// SaveCharacter saves a character to the database
func SaveCharacter(db *DB, character models.Character, discordID string) (models.Character, error) {
	if db == nil {
		return models.Character{}, fmt.Errorf("database connection is nil")
	}

	if discordID == "" {
		return models.Character{}, fmt.Errorf("discord ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Associate character with owner
	character.Owner = discordID

	// Create an initial weapon for the character
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	initialWeapon := roller.GenerateInitialWeaponItem(character.Characteristics.Alignment.Trait_Name, rng)

	// Generate a unique inventory key for the weapon
	inventoryKey := fmt.Sprintf("weapon_%d", time.Now().UnixNano())

	// Save the item to the database
	err := SaveItem(db, initialWeapon, inventoryKey, discordID)
	if err != nil {
		return models.Character{}, fmt.Errorf("failed to save initial weapon: %w", err)
	}

	// Set the character's equipped weapon
	character.EquippedWeapon = models.EquippedItem{
		ItemKey:  inventoryKey,
		ItemName: initialWeapon.Name,
	}

	// Initialize user if they don't exist
	user, err := GetUserByID(db, discordID)
	if err != nil {
		// Create new user
		user, err = CreateUser(db, discordID)
		if err != nil {
			return models.Character{}, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Initialize user inventory if it doesn't exist
	if user.Inventory == nil {
		user.Inventory = make(map[string]string)
	}

	// Add weapon to user's inventory
	user.Inventory[inventoryKey] = initialWeapon.Name

	// Update user in database
	userCollection := db.GetCollection(usersCollection)
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"discordID": discordID},
		bson.M{"$set": bson.M{"inventory": user.Inventory}},
	)
	if err != nil {
		return models.Character{}, fmt.Errorf("failed to update user inventory: %w", err)
	}

	// Insert character
	collection := db.GetCollection(charactersCollection)
	result, err := collection.InsertOne(ctx, character)
	if err != nil {
		return models.Character{}, fmt.Errorf("failed to save character: %w", err)
	}

	// Get the ID of the inserted document
	character.ID = result.InsertedID.(primitive.ObjectID)

	// Also update the user document to reference this character
	userCollection = db.GetCollection(usersCollection)
	filter := bson.M{"discordID": discordID}
	update := bson.M{"$set": bson.M{"character": character}}
	opts := options.Update().SetUpsert(true)

	_, err = userCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return models.Character{}, fmt.Errorf("failed to update user with character: %w", err)
	}

	return character, nil
}

// DeleteCharacter removes a character from the database
func DeleteCharacter(db *DB, userID string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First, get the character to ensure it exists
	character, err := GetCharacterByOwner(db, userID)
	if err != nil {
		return fmt.Errorf("no character found for this user: %w", err)
	}

	// Delete the character from the characters collection
	charCollection := db.GetCollection(charactersCollection)
	_, err = charCollection.DeleteOne(ctx, character)
	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	// Also remove the character reference from the user document
	userCollection := db.GetCollection(usersCollection)
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"discordID": userID},
		bson.M{"$unset": bson.M{"character": ""}},
	)
	if err != nil {
		return fmt.Errorf("failed to update user after character deletion: %w", err)
	}

	return nil
}

// GetCharacterByOwner retrieves a character by owner ID with equipment stats applied
func GetCharacterByOwner(db *DB, ownerID string) (models.Character, error) {
	if db == nil {
		return models.Character{}, fmt.Errorf("database connection is nil")
	}

	if ownerID == "" {
		return models.Character{}, fmt.Errorf("owner ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(charactersCollection)
	filter := bson.M{"owner": ownerID}

	var character models.Character
	err := collection.FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Character{}, fmt.Errorf("no character found for user: %s", ownerID)
		}
		return models.Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	// Clear any existing bonuses
	character = clearEquipmentBonuses(character)

	// Apply trait bonuses first
	character = applyTraitBonuses(character)

	// Apply equipment bonuses if there's an equipped item
	if character.EquippedWeapon.ItemKey != "" {
		character = applyEquipmentBonuses(db, character)
	}

	return character, nil
}

// GetCharacter retrieves a character by ID with equipment stats applied
func GetCharacter(db *DB, characterID string) (models.Character, error) {
	if db == nil {
		return models.Character{}, fmt.Errorf("database connection is nil")
	}

	if characterID == "" {
		return models.Character{}, fmt.Errorf("character ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(characterID)
	if err != nil {
		return models.Character{}, fmt.Errorf("invalid character ID format: %w", err)
	}

	collection := db.GetCollection(charactersCollection)
	filter := bson.M{"_id": objectID}

	var character models.Character
	err = collection.FindOne(ctx, filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Character{}, fmt.Errorf("character not found")
		}
		return models.Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	// Clear any existing bonuses
	character = clearEquipmentBonuses(character)

	// Apply trait bonuses first
	character = applyTraitBonuses(character)

	// Apply equipment bonuses if there's an equipped item
	if character.EquippedWeapon.ItemKey != "" {
		character = applyEquipmentBonuses(db, character)
	}

	return character, nil
}

// ResetRerollCounts resets a user's reroll counts to daily limit
func ResetRerollCounts(db *DB, userID string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := db.GetCollection(usersCollection)

	// Update the user's reroll counts
	_, err := userCollection.UpdateOne(
		ctx,
		bson.M{"discordID": userID},
		bson.M{"$set": bson.M{
			"fullRerolls":     2,
			"statRerolls":     1,
			"lastRerollReset": time.Now(),
		}},
	)

	if err != nil {
		return fmt.Errorf("failed to reset reroll counts: %w", err)
	}

	return nil
}

// UseFullReroll decrements a user's full reroll count
func UseFullReroll(db *DB, userID string) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	// Get the user to check current reroll count
	user, err := GetUserByID(db, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	if user.FullRerolls <= 0 {
		return 0, fmt.Errorf("no full rerolls remaining today")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := db.GetCollection(usersCollection)

	// Decrement full reroll count
	remainingRerolls := user.FullRerolls - 1
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"discordID": userID},
		bson.M{"$set": bson.M{"fullRerolls": remainingRerolls}},
	)

	if err != nil {
		return user.FullRerolls, fmt.Errorf("failed to use full reroll: %w", err)
	}

	return remainingRerolls, nil
}

// UseStatReroll decrements a user's stat reroll count
func UseStatReroll(db *DB, userID string) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	// Get the user to check current reroll count
	user, err := GetUserByID(db, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	if user.StatRerolls <= 0 {
		return 0, fmt.Errorf("no stat rerolls remaining today")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := db.GetCollection(usersCollection)

	// Decrement stat reroll count
	remainingRerolls := user.StatRerolls - 1
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"discordID": userID},
		bson.M{"$set": bson.M{"statRerolls": remainingRerolls}},
	)

	if err != nil {
		return user.StatRerolls, fmt.Errorf("failed to use stat reroll: %w", err)
	}

	return remainingRerolls, nil
}

// RerollSingleStat rerolls a specific stat for a character
func RerollSingleStat(db *DB, userID string, statType variables.StatType) (models.Stat, error) {
	// Create RNG for reroll
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate the new stat based on type
	var newStat models.Stat
	var statField string

	switch statType {
	case variables.Vitality:
		newStat = roller.GenerateStat(variables.Vitality, roller.VitalityRarity, rng)
		statField = "stats.vitality"
	case variables.Durability:
		newStat = roller.GenerateStat(variables.Durability, roller.DurabilityRarity, rng)
		statField = "stats.durability"
	case variables.Speed:
		newStat = roller.GenerateStat(variables.Speed, roller.SpeedRarity, rng)
		statField = "stats.speed"
	case variables.Strength:
		newStat = roller.GenerateStat(variables.Strength, roller.StrengthRarity, rng)
		statField = "stats.strength"
	case variables.Intelligence:
		newStat = roller.GenerateStat(variables.Intelligence, roller.IntelligenceRarity, rng)
		statField = "stats.intelligence"
	case variables.Mana:
		newStat = roller.GenerateStat(variables.Mana, roller.ManaRarity, rng)
		statField = "stats.mana"
	case variables.Mastery:
		newStat = roller.GenerateStat(variables.Mastery, roller.MasteryRarity, rng)
		statField = "stats.mastery"
	default:
		return models.Stat{}, fmt.Errorf("invalid stat type")
	}

	// Update the character in the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	charCollection := db.GetCollection(charactersCollection)
	_, err := charCollection.UpdateOne(
		ctx,
		bson.M{"owner": userID},
		bson.M{"$set": bson.M{statField: newStat}},
	)

	if err != nil {
		return models.Stat{}, fmt.Errorf("failed to update character: %w", err)
	}

	return newStat, nil
}

// AddExperience adds XP to a character and handles level ups
func AddExperience(db *DB, userID string, expAmount int) (int, int, bool, error) {
	if db == nil {
		return 0, 0, false, fmt.Errorf("database connection is nil")
	}

	// Get the character
	character, err := GetCharacterByOwner(db, userID)
	if err != nil {
		return 0, 0, false, fmt.Errorf("no character found for this user: %w", err)
	}

	// Calculate new experience
	newExp := character.Experience + expAmount
	oldLevel := character.Level

	// Calculate new level based on XP
	newLevel := calculateLevel(newExp)

	// Check if level up occurred
	leveledUp := newLevel > oldLevel

	// Update the character with new experience and level
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	charCollection := db.GetCollection(charactersCollection)
	_, err = charCollection.UpdateOne(
		ctx,
		bson.M{"owner": userID},
		bson.M{"$set": bson.M{
			"experience": newExp,
			"level":      newLevel,
		}},
	)

	if err != nil {
		return character.Experience, character.Level, false, fmt.Errorf("failed to update experience: %w", err)
	}

	return newExp, newLevel, leveledUp, nil
}

// calculateLevel determines level based on experience points
func calculateLevel(exp int) int {
	level := 1 // Start at level 1

	// Required XP for each level increases using a multiplier
	requiredXP := variables.LevelUpBaseXP

	// Calculate level based on XP
	for exp >= requiredXP {
		level++
		// Next level requires more XP
		requiredXP = int(float64(requiredXP) * variables.LevelUpMultiplier)
	}

	return level
}

// GetXPForNextLevel calculates how much XP is needed for the next level
func GetXPForNextLevel(currentLevel int) int {
	requiredXP := variables.LevelUpBaseXP

	// Calculate required XP for next level
	for i := 1; i < currentLevel; i++ {
		requiredXP = int(float64(requiredXP) * variables.LevelUpMultiplier)
	}

	return requiredXP
}
