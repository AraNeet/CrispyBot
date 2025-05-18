package database

import (
	"CrispyBot/variables"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	// PostgreSQL driver
)

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, userID string) (User, error) {
	if userID == "" {
		return User{}, fmt.Errorf("discord ID isn't passed")
	}

	// Check if user already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE discord_id = $1)", userID).Scan(&exists)
	if err != nil {
		return User{}, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if !exists {
		// Insert the new user
		_, err = db.Exec("INSERT INTO users (discord_id) VALUES ($1)", userID)
		if err != nil {
			return User{}, fmt.Errorf("failed to create user: %w", err)
		}
	}

	return User{DiscordID: userID}, nil
}

// GetUserByID retrieves a user by Discord ID
func GetUserByID(db *sql.DB, discordID string) (User, error) {
	if discordID == "" {
		return User{}, fmt.Errorf("discord ID is required")
	}

	// Query the database for the user
	var user User
	err := db.QueryRow("SELECT discord_id FROM users WHERE discord_id = $1", discordID).Scan(&user.DiscordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("failed to find user: %w", err)
	}

	// Try to get the character for this user
	character, err := GetCharacterByOwner(db, discordID)
	if err == nil {
		user.Character = &character
	}

	return user, nil
}

// SaveCharacter saves a character to the database
func SaveCharacter(db *sql.DB, character Character, discordID string) (Character, error) {
	if discordID == "" {
		return Character{}, fmt.Errorf("discord ID is required")
	}

	// Generate a UUID for the character if not provided
	if character.ID == "" {
		character.ID = uuid.New().String()
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return Character{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if the user exists, create if not
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE discord_id = $1)", discordID).Scan(&exists)
	if err != nil {
		return Character{}, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if !exists {
		_, err = tx.Exec("INSERT INTO users (discord_id) VALUES ($1)", discordID)
		if err != nil {
			return Character{}, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Insert the character
	_, err = tx.Exec("INSERT INTO characters (id, owner) VALUES ($1, $2)", character.ID, discordID)
	if err != nil {
		return Character{}, fmt.Errorf("failed to save character: %w", err)
	}

	// Save stats
	err = saveStats(tx, character.ID, character.Stats)
	if err != nil {
		return Character{}, fmt.Errorf("failed to save stats: %w", err)
	}

	// Save attributes
	err = saveAttributes(tx, character.ID, character.Attributes)
	if err != nil {
		return Character{}, fmt.Errorf("failed to save attributes: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return Character{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the saved character with the updated ID
	character.Owner = discordID
	return character, nil
}

// saveStats saves all character stats to the database
func saveStats(tx *sql.Tx, characterID string, stats StatsSheets) error {
	// Define the stats to be saved
	statsToSave := []struct {
		StatType variables.StatType
		Stat     Stat
	}{
		{variables.Vitality, stats.Vitality},
		{variables.Durability, stats.Durability},
		{variables.Speed, stats.Speed},
		{variables.Strength, stats.Strength},
		{variables.Intelligence, stats.Intelligence},
		{variables.ManaFlow, stats.ManaFlow},
		{variables.SkillLevel, stats.SkillLevel},
	}

	// Insert each stat - PostgreSQL uses $n for parameter placeholders
	for _, s := range statsToSave {
		_, err := tx.Exec(
			"INSERT INTO stats (character_id, stat_type, stat_name, rarity, value) VALUES ($1, $2, $3, $4, $5)",
			characterID, s.StatType, s.Stat.Stat_Name, s.Stat.Rarity, s.Stat.Value,
		)
		if err != nil {
			return fmt.Errorf("failed to save stat: %w", err)
		}
	}

	return nil
}

// saveAttributes saves all character attributes to the database
func saveAttributes(tx *sql.Tx, characterID string, attrs Attributes) error {
	// Define the traits to be saved
	traitsToSave := []struct {
		Category string
		Trait    Trait
	}{
		{"race", attrs.Race},
		{"element", attrs.Element},
		{"trait", attrs.Trait},
		{"weakness", attrs.Weakness},
		{"alignment", attrs.Alignment},
		{"x_factor", attrs.X_Factor},
	}

	// Insert each trait
	for _, t := range traitsToSave {
		// Insert the trait - PostgreSQL uses RETURNING to get the inserted ID
		var traitID int64
		err := tx.QueryRow(
			"INSERT INTO traits (character_id, trait_type, trait_name, rarity, trait_category) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			characterID, t.Trait.Type, t.Trait.Trait_Name, t.Trait.Rarity, t.Category,
		).Scan(&traitID)
		if err != nil {
			return fmt.Errorf("failed to save trait: %w", err)
		}

		// Insert trait stat values
		for statName, value := range t.Trait.Stats_Value {
			_, err := tx.Exec(
				"INSERT INTO trait_stat_values (trait_id, stat_name, value) VALUES ($1, $2, $3)",
				traitID, statName, value,
			)
			if err != nil {
				return fmt.Errorf("failed to save trait stat value: %w", err)
			}
		}
	}

	return nil
}

// GetCharacterByOwner retrieves a character by owner ID
func GetCharacterByOwner(db *sql.DB, ownerID string) (Character, error) {
	if ownerID == "" {
		return Character{}, fmt.Errorf("owner ID is required")
	}

	// Query for the character
	var character Character
	err := db.QueryRow("SELECT id, owner FROM characters WHERE owner = $1 LIMIT 1", ownerID).Scan(&character.ID, &character.Owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return Character{}, fmt.Errorf("no character found for user: %s", ownerID)
		}
		return Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	// Load stats
	stats, err := loadStats(db, character.ID)
	if err != nil {
		return Character{}, fmt.Errorf("failed to load stats: %w", err)
	}
	character.Stats = stats

	// Load attributes
	attributes, err := loadAttributes(db, character.ID)
	if err != nil {
		return Character{}, fmt.Errorf("failed to load attributes: %w", err)
	}
	character.Attributes = attributes

	return character, nil
}

// loadStats loads all stats for a character
func loadStats(db *sql.DB, characterID string) (StatsSheets, error) {
	var stats StatsSheets

	// Query all stats for this character
	rows, err := db.Query(
		"SELECT stat_type, stat_name, rarity, value FROM stats WHERE character_id = $1",
		characterID,
	)
	if err != nil {
		return stats, fmt.Errorf("failed to query stats: %w", err)
	}
	defer rows.Close()

	// Process each stat
	for rows.Next() {
		var statType int
		var stat Stat
		err := rows.Scan(&statType, &stat.Stat_Name, &stat.Rarity, &stat.Value)
		if err != nil {
			return stats, fmt.Errorf("failed to scan stat: %w", err)
		}

		// Assign the stat to the correct field based on type
		stat.Type = variables.StatType(statType)
		switch stat.Type {
		case variables.Vitality:
			stats.Vitality = stat
		case variables.Durability:
			stats.Durability = stat
		case variables.Speed:
			stats.Speed = stat
		case variables.Strength:
			stats.Strength = stat
		case variables.Intelligence:
			stats.Intelligence = stat
		case variables.ManaFlow:
			stats.ManaFlow = stat
		case variables.SkillLevel:
			stats.SkillLevel = stat
		}
	}

	if err = rows.Err(); err != nil {
		return stats, fmt.Errorf("error iterating stats rows: %w", err)
	}

	return stats, nil
}

// loadAttributes loads all attributes for a character
func loadAttributes(db *sql.DB, characterID string) (Attributes, error) {
	var attributes Attributes

	// Query all traits for this character
	rows, err := db.Query(
		"SELECT id, trait_type, trait_name, rarity, trait_category FROM traits WHERE character_id = $1",
		characterID,
	)
	if err != nil {
		return attributes, fmt.Errorf("failed to query traits: %w", err)
	}
	defer rows.Close()

	// Process each trait
	for rows.Next() {
		var trait Trait
		var category string
		err := rows.Scan(&trait.ID, &trait.Type, &trait.Trait_Name, &trait.Rarity, &category)
		if err != nil {
			return attributes, fmt.Errorf("failed to scan trait: %w", err)
		}

		// Load trait stat values
		trait.Stats_Value, err = loadTraitStatValues(db, trait.ID)
		if err != nil {
			return attributes, fmt.Errorf("failed to load trait stat values: %w", err)
		}

		// Assign the trait to the correct field based on category
		switch strings.ToLower(category) {
		case "race":
			attributes.Race = trait
		case "element":
			attributes.Element = trait
		case "trait":
			attributes.Trait = trait
		case "weakness":
			attributes.Weakness = trait
		case "alignment":
			attributes.Alignment = trait
		case "x_factor":
			attributes.X_Factor = trait
		}
	}

	if err = rows.Err(); err != nil {
		return attributes, fmt.Errorf("error iterating trait rows: %w", err)
	}

	return attributes, nil
}

// loadTraitStatValues loads all stat values for a trait
func loadTraitStatValues(db *sql.DB, traitID int) (map[string]int, error) {
	statValues := make(map[string]int)

	// Query all stat values for this trait
	rows, err := db.Query(
		"SELECT stat_name, value FROM trait_stat_values WHERE trait_id = $1",
		traitID,
	)
	if err != nil {
		return statValues, fmt.Errorf("failed to query trait stat values: %w", err)
	}
	defer rows.Close()

	// Process each stat value
	for rows.Next() {
		var statName string
		var value int
		err := rows.Scan(&statName, &value)
		if err != nil {
			return statValues, fmt.Errorf("failed to scan trait stat value: %w", err)
		}

		statValues[statName] = value
	}

	if err = rows.Err(); err != nil {
		return statValues, fmt.Errorf("error iterating trait stat value rows: %w", err)
	}

	return statValues, nil
}

// GetCharacter retrieves a character by character ID
func GetCharacter(db *sql.DB, characterID string) (Character, error) {
	if characterID == "" {
		return Character{}, fmt.Errorf("character ID is required")
	}

	// Query for the character
	var character Character
	err := db.QueryRow("SELECT id, owner FROM characters WHERE id = $1", characterID).Scan(&character.ID, &character.Owner)
	if err != nil {
		if err == sql.ErrNoRows {
			return Character{}, fmt.Errorf("character not found: %s", characterID)
		}
		return Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	// Load stats
	stats, err := loadStats(db, character.ID)
	if err != nil {
		return Character{}, fmt.Errorf("failed to load stats: %w", err)
	}
	character.Stats = stats

	// Load attributes
	attributes, err := loadAttributes(db, character.ID)
	if err != nil {
		return Character{}, fmt.Errorf("failed to load attributes: %w", err)
	}
	character.Attributes = attributes

	return character, nil
}
