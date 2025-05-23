package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User Model
/*
	ID - ObjectID for the user record.
	DiscordID - User's Discord ID for identification.
	Wallet - User's current currency/money balance.
	Inventory - User's item storage. Note: Key is inventory slot, value is item identifier.
	Character - User's active character data.
	FullRerolls - Number of complete character rerolls available.
	StatRerolls - Number of stat-only rerolls available.
	LastRerollReset - Timestamp of last reroll counter reset. Note: Used for daily/periodic reroll refresh.
*/
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
