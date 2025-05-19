package database

import (
	"CrispyBot/database/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// EquipItem equips an item to a character
func EquipItem(db *DB, userID string, itemKey string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Get the user
	user, err := GetUserByID(db, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if the item exists in the user's inventory
	itemName, ok := user.Inventory[itemKey]
	if !ok {
		return fmt.Errorf("item not found in inventory")
	}

	// Get the character
	character, err := GetCharacterByOwner(db, userID)
	if err != nil {
		return fmt.Errorf("no character found for this user")
	}

	// Set up context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get current equipped weapon
	// currentWeapon := character.EquippedWeapon

	// Create update operations
	charCollection := db.GetCollection(charactersCollection)

	// Set the equipped weapon and item name
	updates := bson.M{
		"$set": bson.M{
			"equippedWeapon.itemKey":  itemKey,
			"equippedWeapon.itemName": itemName,
		},
	}

	// Update the character
	_, err = charCollection.UpdateOne(
		ctx,
		bson.M{"_id": character.ID},
		updates,
	)

	if err != nil {
		return fmt.Errorf("failed to equip item: %w", err)
	}

	return nil
}

// UnequipItem removes the currently equipped item
func UnequipItem(db *DB, userID string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Get the character
	character, err := GetCharacterByOwner(db, userID)
	if err != nil {
		return fmt.Errorf("no character found for this user")
	}

	// Set up context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create update operations
	charCollection := db.GetCollection(charactersCollection)

	// Clear the equipped weapon
	updates := bson.M{
		"$set": bson.M{
			"equippedWeapon": models.EquippedItem{},
		},
	}

	// Update the character
	_, err = charCollection.UpdateOne(
		ctx,
		bson.M{"_id": character.ID},
		updates,
	)

	if err != nil {
		return fmt.Errorf("failed to unequip item: %w", err)
	}

	return nil
}

// GetItemFromShop retrieves an item from the shop by index
func GetItemFromShop(db *DB, itemIndex int) (models.Item, error) {
	if db == nil {
		return models.Item{}, fmt.Errorf("database connection is nil")
	}

	// Get the shop
	shop, err := GetShop(db)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to get shop: %w", err)
	}

	// Check if item exists
	item, ok := shop.Inventory.Items[itemIndex]
	if !ok {
		return models.Item{}, fmt.Errorf("item not found in shop")
	}

	return item, nil
}

// SaveItem saves item stats when a user purchases it
func SaveItem(db *DB, item models.Item, inventoryKey string, userID string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create or get the items collection
	itemsCollection := db.GetCollection("items")

	// Create item record
	itemRecord := models.ItemRecord{
		OwnerID:      userID,
		InventoryKey: inventoryKey,
		Item:         item,
		Timestamp:    time.Now(),
	}

	// Insert the item
	_, err := itemsCollection.InsertOne(ctx, itemRecord)
	if err != nil {
		return fmt.Errorf("failed to save item stats: %w", err)
	}

	return nil
}

// GetItem retrieves an item's stats by inventory key
func GetItem(db *DB, userID string, inventoryKey string) (models.Item, error) {
	if db == nil {
		return models.Item{}, fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the items collection
	itemsCollection := db.GetCollection("items")

	// Query for the item
	filter := bson.M{
		"ownerID":      userID,
		"inventoryKey": inventoryKey,
	}

	var itemRecord models.ItemRecord
	err := itemsCollection.FindOne(ctx, filter).Decode(&itemRecord)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to get item: %w", err)
	}

	return itemRecord.Item, nil
}
