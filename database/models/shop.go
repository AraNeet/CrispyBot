package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shop struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timer     time.Time          `bson:"timeRemaining" json:"timeRemaining"`
	Inventory Inventory          `bson:"inventory" json:"inventory"`
}

type Inventory struct {
	Items map[int]Item `bson:"items" json:"items"`
}

type Item struct {
	Name   string         `bson:"name" json:"name"`
	Rarity string         `bson:"rarity" json:"rarity"`
	Stats  map[string]int `bson:"stats" json:"stats"`
	Price  int            `bson:"price" json:"price"`
}
