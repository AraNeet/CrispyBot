package database

import (
	"CrispyBot/variables"
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

// SaveCharacter saves a character to the database
func SaveCharacter(db *surrealdb.DB, character Character, discordID string) (Character, error) {
	if discordID == "" {
		return Character{}, fmt.Errorf("discord ID is required")
	}

	// Prepare the data to save
	data := map[string]interface{}{
		"owner":      discordID,
		"stats":      character.Stats,
		"attributes": character.Attributes,
	}

	// Create the record in the database
	query := "CREATE type::table($tb) CONTENT $data RETURN AFTER"
	params := map[string]interface{}{
		"tb":   "characters",
		"data": data,
	}

	var response []map[string]interface{}
	if err := db.Query(query, params); err != nil {
		return Character{}, fmt.Errorf("failed to save character: %w", err)
	}

	if len(response) == 0 || len(response[0]) == 0 {
		return Character{}, fmt.Errorf("failed to create character record")
	}

	// Parse the created record similarly to how we do in GetCharacterByOwner
	characterData := response[0]
	var savedChar Character

	if id, ok := characterData["id"].(string); ok {
		savedChar.ID = id
	}

	// Set the owner explicitly
	savedChar.Owner = discordID

	if statsData, ok := characterData["stats"].(map[string]interface{}); ok {
		savedChar.Stats = convertMapToStats(statsData)
	} else {
		// If stats weren't returned properly, use the original stats
		savedChar.Stats = character.Stats
	}

	if attrsData, ok := characterData["attributes"].(map[string]interface{}); ok {
		savedChar.Attributes = convertMapToAttributes(attrsData)
	} else {
		// If attributes weren't returned properly, use the original attributes
		savedChar.Attributes = character.Attributes
	}

	return savedChar, nil
}

// Update this function in database/repository.go
func GetCharacterByOwner(db *surrealdb.DB, ownerID string) (Character, error) {
	if ownerID == "" {
		return Character{}, fmt.Errorf("owner ID is required")
	}

	// Use parameterized query with type::table function
	query := "SELECT * FROM type::table($tb) WHERE owner = $owner LIMIT 1"
	params := map[string]interface{}{
		"tb":    "characters",
		"owner": ownerID,
	}

	var response []map[string]interface{}
	if err := db.Query(query, params); err != nil {
		return Character{}, fmt.Errorf("failed to query character: %w", err)
	}

	if len(response) == 0 || len(response[0]) == 0 {
		return Character{}, fmt.Errorf("no character found for user: %s", ownerID)
	}

	// Parse the result into a Character object
	// Note: SurrealDB's response structure might require accessing a 'result' field first
	characterData := response[0]

	var character Character

	// Extract the basic fields
	if id, ok := characterData["id"].(string); ok {
		character.ID = id
	}

	if owner, ok := characterData["owner"].(string); ok {
		character.Owner = owner
	}

	// For nested structures like Stats and Attributes, we'll need to handle the conversion
	if statsData, ok := characterData["stats"].(map[string]interface{}); ok {
		// Convert the map to a StatsSheets struct
		character.Stats = convertMapToStats(statsData)
	}

	if attrsData, ok := characterData["attributes"].(map[string]interface{}); ok {
		// Convert the map to an Attributes struct
		character.Attributes = convertMapToAttributes(attrsData)
	}

	return character, nil
}

// Helper function to convert a map to StatsSheets
func convertMapToStats(statsMap map[string]interface{}) StatsSheets {
	stats := StatsSheets{}

	// Handle each stat, with appropriate type checking and conversion
	if vitalityMap, ok := statsMap["vitality"].(map[string]interface{}); ok {
		stats.Vitality = convertMapToStat(vitalityMap)
	}

	if durabilityMap, ok := statsMap["durability"].(map[string]interface{}); ok {
		stats.Durability = convertMapToStat(durabilityMap)
	}

	if speedMap, ok := statsMap["speed"].(map[string]interface{}); ok {
		stats.Speed = convertMapToStat(speedMap)
	}

	if strengthMap, ok := statsMap["strength"].(map[string]interface{}); ok {
		stats.Strength = convertMapToStat(strengthMap)
	}

	if intelligenceMap, ok := statsMap["intelligence"].(map[string]interface{}); ok {
		stats.Intelligence = convertMapToStat(intelligenceMap)
	}

	if manaflowMap, ok := statsMap["manaflow"].(map[string]interface{}); ok {
		stats.ManaFlow = convertMapToStat(manaflowMap)
	}

	if skilllevelMap, ok := statsMap["skilllevel"].(map[string]interface{}); ok {
		stats.SkillLevel = convertMapToStat(skilllevelMap)
	}

	return stats
}

// Helper function to convert a map to a Stat
func convertMapToStat(statMap map[string]interface{}) Stat {
	stat := Stat{}

	if rarity, ok := statMap["rarity"].(string); ok {
		stat.Rarity = rarity
	}

	if statName, ok := statMap["stat_name"].(string); ok {
		stat.Stat_Name = statName
	}

	if typeVal, ok := statMap["type"].(float64); ok {
		stat.Type = variables.StatType(int(typeVal))
	}

	if value, ok := statMap["value"].(float64); ok {
		stat.Value = int(value)
	}

	return stat
}

// Helper function to convert a map to Attributes
func convertMapToAttributes(attrsMap map[string]interface{}) Attributes {
	attrs := Attributes{}

	if raceMap, ok := attrsMap["race"].(map[string]interface{}); ok {
		attrs.Race = convertMapToTrait(raceMap)
	}

	if elementMap, ok := attrsMap["element"].(map[string]interface{}); ok {
		attrs.Element = convertMapToTrait(elementMap)
	}

	if traitMap, ok := attrsMap["trait"].(map[string]interface{}); ok {
		attrs.Trait = convertMapToTrait(traitMap)
	}

	if weaknessMap, ok := attrsMap["weakness"].(map[string]interface{}); ok {
		attrs.Weakness = convertMapToTrait(weaknessMap)
	}

	if alignmentMap, ok := attrsMap["alignment"].(map[string]interface{}); ok {
		attrs.Alignment = convertMapToTrait(alignmentMap)
	}

	if xFactorMap, ok := attrsMap["x_factor"].(map[string]interface{}); ok {
		attrs.X_Factor = convertMapToTrait(xFactorMap)
	}

	return attrs
}

// Helper function to convert a map to a Trait
func convertMapToTrait(traitMap map[string]interface{}) Trait {
	trait := Trait{}

	if rarity, ok := traitMap["rarity"].(string); ok {
		trait.Rarity = rarity
	}

	if traitName, ok := traitMap["trait_name"].(string); ok {
		trait.Trait_Name = traitName
	}

	if typeVal, ok := traitMap["type"].(float64); ok {
		trait.Type = variables.TraitType(int(typeVal))
	}

	// Handle the stats_value map
	trait.Stats_Value = make(map[string]int)
	if statsValueMap, ok := traitMap["stats_value"].(map[string]interface{}); ok {
		for key, val := range statsValueMap {
			if floatVal, ok := val.(float64); ok {
				trait.Stats_Value[key] = int(floatVal)
			}
		}
	}

	return trait
}
