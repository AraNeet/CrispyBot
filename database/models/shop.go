package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Shop Model
/*
	ID - ObjectID for the shop instance.
	Timer - Time when shop inventory refreshes/resets.
	Inventory - Current items available for purchase.
*/
type Shop struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timer     time.Time          `bson:"timeRemaining" json:"timeRemaining"`
	Inventory Inventory          `bson:"inventory" json:"inventory"`
}

// Inventory Model
/*
	Items - Map of available items. Note: Key is slot number, value is the item.
*/
type Inventory struct {
	Items map[int]Item `bson:"items" json:"items"`
}

// Item Model
/*
	Name - Display name of the item.
	Rarity - Rarity level of the item (common, rare, epic, etc.).
	Stats - Stat modifications provided by item. Note: Key is stat name, value is modifier amount.
	Price - Cost to purchase this item.
*/
type Item struct {
	Name   string         `bson:"name" json:"name"`
	Rarity string         `bson:"rarity" json:"rarity"`
	Stats  map[string]int `bson:"stats" json:"stats"`
	Price  int            `bson:"price" json:"price"`
}
