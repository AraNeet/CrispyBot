package database

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	DisID     string    `json:"DiscordID"`
	Character Character `json:"Character"`
}

type Character struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Owner      string    `json:"owner,omitempty"`
	Stats      string    `json:"stats"`
	Attributes string    `json:"attributes"`
}

type Stats struct {
	ID           string `json:"id,omitempty"`
	Vitality     string `json:"vitality"`
	Durability   string `json:"durability"`
	Speed        string `json:"speed"`
	Strength     string `json:"strength"`
	Intelligence string `json:"intelligence"`
	ManaFlow     string `json:"manaflow"`
	SkillLevel   string `json:"skilllevel"`
}

type Attributes struct {
	ID        string `json:"id,omitempty"`
	Race      string `json:"race"`
	Element   string `json:"element"`
	Buff      string `json:"buff"`
	Weakness  string `json:"weakness"`
	Alignment string `json:"alignment"`
	X_Factor  string `json:"x-factor"`
}
