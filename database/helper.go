package database

import (
	"CrispyBot/database/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// resetAllUsersRerolls resets reroll counts for all users
func resetAllUsersRerolls(db *DB) {
	if db == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := db.GetCollection(usersCollection)

	// Update all users
	_, err := userCollection.UpdateMany(
		ctx,
		bson.M{}, // Match all documents
		bson.M{"$set": bson.M{
			"fullRerolls":     2,
			"statRerolls":     1,
			"lastRerollReset": time.Now(),
		}},
	)

	if err != nil {
		fmt.Printf("Error resetting reroll counts: %v\n", err)
	} else {
		fmt.Println("Successfully reset reroll counts for all users")
	}
}

// StartShopRefreshScheduler starts a goroutine to check and refresh the shop periodically
func StartShopRefreshScheduler(db *DB) {
	go func() {
		for {
			// Check and update the shop if needed
			err := checkAndRefreshShop(db)
			if err != nil {
				fmt.Printf("Error refreshing shop: %v\n", err)
			}

			// Calculate time until next check
			// We'll check every 5 minutes to see if the shop needs refreshing
			nextCheck := time.Minute * 5

			// Sleep until next check
			time.Sleep(nextCheck)
		}
	}()

	fmt.Println("Shop refresh scheduler started")
}

// checkAndRefreshShop checks if the shop needs to be refreshed and updates it
func checkAndRefreshShop(db *DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Get the current shop
	shop, err := GetShop(db)
	if err != nil {
		return fmt.Errorf("failed to get shop: %w", err)
	}

	// Check if the shop timer has expired
	if time.Now().After(shop.Timer) {
		fmt.Println("Shop timer expired, refreshing shop")
		RefreshShop(db, shop)

		// Reset all users' reroll counts
		resetAllUsersRerolls(db)
	}

	return nil
}

// applyEquipmentBonuses adds equipment stats to character stats
func applyEquipmentBonuses(db *DB, character models.Character) models.Character {
	// If no weapon is equipped, return character as is
	if character.EquippedWeapon.ItemKey == "" {
		return character
	}

	// Get the equipped item's stats
	item, err := GetItem(db, character.Owner, character.EquippedWeapon.ItemKey)
	if err != nil {
		fmt.Printf("Error getting equipped item stats: %v\n", err)
		return character
	}

	// Apply stat bonuses
	for statName, value := range item.Stats {
		switch statName {
		case "Vitality":
			character.Stats.Vitality.EquipBonus = value
		case "Strength":
			character.Stats.Strength.EquipBonus = value
		case "Speed":
			character.Stats.Speed.EquipBonus = value
		case "Durability":
			character.Stats.Durability.EquipBonus = value
		case "Intelligence":
			character.Stats.Intelligence.EquipBonus = value
		case "Mana":
			character.Stats.Mana.EquipBonus = value
		case "Mastery":
			character.Stats.Mastery.EquipBonus = value
		}
	}

	// Calculate total values
	character.Stats.Vitality.TotalValue = character.Stats.Vitality.Value + character.Stats.Vitality.EquipBonus
	character.Stats.Strength.TotalValue = character.Stats.Strength.Value + character.Stats.Strength.EquipBonus
	character.Stats.Speed.TotalValue = character.Stats.Speed.Value + character.Stats.Speed.EquipBonus
	character.Stats.Durability.TotalValue = character.Stats.Durability.Value + character.Stats.Durability.EquipBonus
	character.Stats.Intelligence.TotalValue = character.Stats.Intelligence.Value + character.Stats.Intelligence.EquipBonus
	character.Stats.Mana.TotalValue = character.Stats.Mana.Value + character.Stats.Mana.EquipBonus
	character.Stats.Mastery.TotalValue = character.Stats.Mastery.Value + character.Stats.Mastery.EquipBonus

	return character
}

// clearEquipmentBonuses removes all bonuses from character stats
func clearEquipmentBonuses(character models.Character) models.Character {
	// Reset all equipment and trait bonuses to 0
	character.Stats.Vitality.EquipBonus = 0
	character.Stats.Strength.EquipBonus = 0
	character.Stats.Speed.EquipBonus = 0
	character.Stats.Durability.EquipBonus = 0
	character.Stats.Intelligence.EquipBonus = 0
	character.Stats.Mana.EquipBonus = 0
	character.Stats.Mastery.EquipBonus = 0

	character.Stats.Vitality.TraitBonus = 0
	character.Stats.Strength.TraitBonus = 0
	character.Stats.Speed.TraitBonus = 0
	character.Stats.Durability.TraitBonus = 0
	character.Stats.Intelligence.TraitBonus = 0
	character.Stats.Mana.TraitBonus = 0
	character.Stats.Mastery.TraitBonus = 0

	// Set total values equal to base values
	character.Stats.Vitality.TotalValue = character.Stats.Vitality.Value
	character.Stats.Strength.TotalValue = character.Stats.Strength.Value
	character.Stats.Speed.TotalValue = character.Stats.Speed.Value
	character.Stats.Durability.TotalValue = character.Stats.Durability.Value
	character.Stats.Intelligence.TotalValue = character.Stats.Intelligence.Value
	character.Stats.Mana.TotalValue = character.Stats.Mana.Value
	character.Stats.Mastery.TotalValue = character.Stats.Mastery.Value

	return character
}

// applyTraitBonuses applies trait bonuses (innate and inadequacy) to character stats
func applyTraitBonuses(character models.Character) models.Character {
	// Apply Innate trait bonuses
	if character.Attributes.Innate.Stats_Value != nil {
		for statName, value := range character.Attributes.Innate.Stats_Value {
			switch statName {
			case "Vitality":
				character.Stats.Vitality.TraitBonus += value
			case "Strength":
				character.Stats.Strength.TraitBonus += value
			case "Speed":
				character.Stats.Speed.TraitBonus += value
			case "Durability":
				character.Stats.Durability.TraitBonus += value
			case "Intelligence":
				character.Stats.Intelligence.TraitBonus += value
			case "Mana":
				character.Stats.Mana.TraitBonus += value
			case "Mastery":
				character.Stats.Mastery.TraitBonus += value
			}
		}
	}

	// Apply Inadequacy trait effects (these are usually negative)
	if character.Attributes.Inadequacy.Stats_Value != nil {
		for statName, value := range character.Attributes.Inadequacy.Stats_Value {
			switch statName {
			case "Vitality":
				character.Stats.Vitality.TraitBonus += value
			case "Strength":
				character.Stats.Strength.TraitBonus += value
			case "Speed":
				character.Stats.Speed.TraitBonus += value
			case "Durability":
				character.Stats.Durability.TraitBonus += value
			case "Intelligence":
				character.Stats.Intelligence.TraitBonus += value
			case "Mana":
				character.Stats.Mana.TraitBonus += value
			case "Mastery":
				character.Stats.Mastery.TraitBonus += value
			}
		}
	}

	// Apply Race characteristic bonuses/penalties
	if character.Characteristics.Race.Stats_Value != nil {
		for statName, value := range character.Characteristics.Race.Stats_Value {
			switch statName {
			case "Vitality":
				character.Stats.Vitality.TraitBonus += value
			case "Strength":
				character.Stats.Strength.TraitBonus += value
			case "Speed":
				character.Stats.Speed.TraitBonus += value
			case "Durability":
				character.Stats.Durability.TraitBonus += value
			case "Intelligence":
				character.Stats.Intelligence.TraitBonus += value
			case "Mana":
				character.Stats.Mana.TraitBonus += value
			case "Mastery":
				character.Stats.Mastery.TraitBonus += value
			}
		}
	}

	// Calculate total values including trait bonuses
	character.Stats.Vitality.TotalValue = character.Stats.Vitality.Value + character.Stats.Vitality.EquipBonus + character.Stats.Vitality.TraitBonus
	character.Stats.Strength.TotalValue = character.Stats.Strength.Value + character.Stats.Strength.EquipBonus + character.Stats.Strength.TraitBonus
	character.Stats.Speed.TotalValue = character.Stats.Speed.Value + character.Stats.Speed.EquipBonus + character.Stats.Speed.TraitBonus
	character.Stats.Durability.TotalValue = character.Stats.Durability.Value + character.Stats.Durability.EquipBonus + character.Stats.Durability.TraitBonus
	character.Stats.Intelligence.TotalValue = character.Stats.Intelligence.Value + character.Stats.Intelligence.EquipBonus + character.Stats.Intelligence.TraitBonus
	character.Stats.Mana.TotalValue = character.Stats.Mana.Value + character.Stats.Mana.EquipBonus + character.Stats.Mana.TraitBonus
	character.Stats.Mastery.TotalValue = character.Stats.Mastery.Value + character.Stats.Mastery.EquipBonus + character.Stats.Mastery.TraitBonus

	return character
}
