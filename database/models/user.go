package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DiscordID string             `bson:"discordID" json:"discordID"`
	Character Character          `bson:"character,omitempty" json:"character,omitempty"`
}
