package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DiscordID       string             `bson:"discordID" json:"discordID"`
	Wallet          int                `bson:"wallet" json:"wallet"`
	Inventory       map[string]string  `bson:"inventory" json:"invertory"`
	Character       Character          `bson:"character,omitempty" json:"character,omitempty"`
	FullRerolls     int                `bson:"fullRerolls" json:"fullRerolls"`
	StatRerolls     int                `bson:"statRerolls" json:"statRerolls"`
	LastRerollReset time.Time          `bson:"lastRerollReset" json:"lastRerollReset"`
}
