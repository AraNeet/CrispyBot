package database

import (
	"CrispyBot/database/models"
	"CrispyBot/shop"

	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	shopCollection = "shop"
)

// GetShop retrieves the current shop or creates a new one if it doesn't exist
func GetShop(db *DB) (models.Shop, error) {
	if db == nil {
		return models.Shop{}, fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(shopCollection)

	// Try to get the current shop
	var shop models.Shop
	err := collection.FindOne(ctx, bson.M{}).Decode(&shop)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No shop exists, create a new one
			return CreateNewShop(db)
		}
		return models.Shop{}, fmt.Errorf("failed to query shop: %w", err)
	}

	// Check if shop needs refreshing
	if shop.Timer.Before(time.Now()) {
		// Shop has expired, refresh it
		shop = RefreshShop(db, shop)
	}

	return shop, nil
}

// CreateNewShop creates a new shop in the database
func CreateNewShop(db *DB) (models.Shop, error) {
	if db == nil {
		return models.Shop{}, fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(shopCollection)

	// Create a new shop
	newShop := shop.CreateShop()

	// Insert into database
	result, err := collection.InsertOne(ctx, newShop)
	if err != nil {
		return models.Shop{}, fmt.Errorf("failed to create shop: %w", err)
	}

	// Get the ID from the inserted document
	id := result.InsertedID
	newShop.ID = id.(primitive.ObjectID)

	return newShop, nil
}

// RefreshShop updates the shop with new inventory
func RefreshShop(db *DB, oldShop models.Shop) models.Shop {
	if db == nil {
		return oldShop
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(shopCollection)

	// Create refreshed shop
	shop.RefreshShop(&oldShop)

	// Update in database
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": oldShop.ID},
		bson.M{"$set": oldShop},
	)

	if err != nil {
		fmt.Printf("Error updating shop: %v\n", err)
		return oldShop
	}

	return oldShop
}

// UpdateShopIfExpired checks if the shop needs to be refreshed and does so if needed
func UpdateShopIfExpired(db *DB) {
	shop, err := GetShop(db)
	if err != nil {
		fmt.Printf("Error getting shop: %v\n", err)
		return
	}

	if shop.Timer.Before(time.Now()) {
		RefreshShop(db, shop)
	}
}

// BuyItem handles the purchase of an item by a user
func BuyItem(db *DB, userID string, itemIndex int) (models.Item, error) {
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

	// Get the user
	user, err := GetUserByID(db, userID)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has enough money
	if user.Wallet < item.Price {
		return models.Item{}, fmt.Errorf("not enough currency to buy this item")
	}

	// Update user's wallet and add item to inventory
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := db.GetCollection(usersCollection)

	// Initialize inventory if it doesn't exist
	if user.Inventory == nil {
		user.Inventory = make(map[string]string)
	}

	// Add item to user's inventory with a unique key
	inventoryKey := fmt.Sprintf("weapon_%d", len(user.Inventory)+1)
	user.Inventory[inventoryKey] = item.Name

	// Deduct price from wallet
	user.Wallet -= item.Price

	// Update user in database
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"wallet":    user.Wallet,
			"inventory": user.Inventory,
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to update user after purchase: %w", err)
	}

	// Save item stats to items collection
	err = SaveItem(db, item, inventoryKey, userID)
	if err != nil {
		fmt.Printf("Warning: Failed to save item stats: %v\n", err)
		// We'll continue even if saving stats fails
	}

	// Remove item from shop
	delete(shop.Inventory.Items, itemIndex)

	// Update shop in database
	shopCollection := db.GetCollection(shopCollection)
	_, err = shopCollection.UpdateOne(
		ctx,
		bson.M{"_id": shop.ID},
		bson.M{"$set": bson.M{"inventory": shop.Inventory}},
	)
	if err != nil {
		fmt.Printf("Error updating shop after purchase: %v\n", err)
	}

	return item, nil
}

// Initialize user wallet if they don't have one
func InitializeUserWallet(db *DB, userID string, initialAmount int) error {
	user, err := GetUserByID(db, userID)
	if err != nil {
		// User doesn't exist, create new user
		_, err = CreateUser(db, userID)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Get the newly created user
		user, err = GetUserByID(db, userID)
		if err != nil {
			return fmt.Errorf("failed to get created user: %w", err)
		}
	}

	// If wallet is 0, set to initial amount
	if user.Wallet == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userCollection := db.GetCollection(usersCollection)

		_, err = userCollection.UpdateOne(
			ctx,
			bson.M{"discordID": userID},
			bson.M{"$set": bson.M{"wallet": initialAmount}},
		)

		if err != nil {
			return fmt.Errorf("failed to initialize wallet: %w", err)
		}
	}

	return nil
}

// Add currency to a user's wallet (for daily rewards, etc.)
func AddCurrency(db *DB, userID string, amount int) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	user, err := GetUserByID(db, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCollection := db.GetCollection(usersCollection)

	// Update wallet
	newBalance := user.Wallet + amount

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"discordID": userID},
		bson.M{"$set": bson.M{"wallet": newBalance}},
	)

	if err != nil {
		return user.Wallet, fmt.Errorf("failed to add currency: %w", err)
	}

	return newBalance, nil
}
