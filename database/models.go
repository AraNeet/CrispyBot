package database

import (
	"CrispyBot/variables"
)

// User model
type User struct {
	DiscordID string     `json:"discord_id" db:"discord_id"`
	Character *Character `json:"character,omitempty" db:"-"`
}

// Character model
type Character struct {
	ID         string      `json:"id,omitempty" db:"id"`
	Owner      string      `json:"owner" db:"owner"`
	Stats      StatsSheets `json:"stats" db:"-"`      // Composite type stored in separate table
	Attributes Attributes  `json:"attributes" db:"-"` // Composite type stored in separate table
}

// StatsSheets model
type StatsSheets struct {
	Vitality     Stat `json:"vitality"`
	Durability   Stat `json:"durability"`
	Speed        Stat `json:"speed"`
	Strength     Stat `json:"strength"`
	Intelligence Stat `json:"intelligence"`
	ManaFlow     Stat `json:"manaflow"`
	SkillLevel   Stat `json:"skilllevel"`
}

// Attributes model
type Attributes struct {
	Race      Trait `json:"race"`
	Element   Trait `json:"element"`
	Trait     Trait `json:"trait"`
	Weakness  Trait `json:"weakness"`
	Alignment Trait `json:"alignment"`
	X_Factor  Trait `json:"x_factor"`
}

// Trait model
type Trait struct {
	ID          int                 `json:"id,omitempty" db:"id"`
	CharacterID string              `json:"character_id,omitempty" db:"character_id"`
	Rarity      string              `json:"rarity" db:"rarity"`
	Trait_Name  string              `json:"trait_name" db:"trait_name"`
	Type        variables.TraitType `json:"type" db:"trait_type"`
	Category    string              `json:"category,omitempty" db:"trait_category"` // Race, Element, etc.
	Stats_Value map[string]int      `json:"stats_value" db:"-"`                     // Stored in separate table
}

// Stat model
type Stat struct {
	ID          int                `json:"id,omitempty" db:"id"`
	CharacterID string             `json:"character_id,omitempty" db:"character_id"`
	Rarity      string             `json:"rarity" db:"rarity"`
	Stat_Name   string             `json:"stat_name" db:"stat_name"`
	Type        variables.StatType `json:"type" db:"stat_type"`
	Value       int                `json:"value" db:"value"`
}

// TraitStatValue model for PostgreSQL relationship
type TraitStatValue struct {
	ID       int    `json:"id,omitempty" db:"id"`
	TraitID  int    `json:"trait_id" db:"trait_id"`
	StatName string `json:"stat_name" db:"stat_name"`
	Value    int    `json:"value" db:"value"`
}

// Database interface models for scanning query results
// These can be useful for directly scanning rows into structs

type CharacterWithOwner struct {
	ID    string `db:"id"`
	Owner string `db:"owner"`
}

type StatWithType struct {
	ID          int    `db:"id"`
	CharacterID string `db:"character_id"`
	StatType    int    `db:"stat_type"`
	StatName    string `db:"stat_name"`
	Rarity      string `db:"rarity"`
	Value       int    `db:"value"`
}

type TraitWithCategory struct {
	ID            int    `db:"id"`
	CharacterID   string `db:"character_id"`
	TraitType     int    `db:"trait_type"`
	TraitName     string `db:"trait_name"`
	Rarity        string `db:"rarity"`
	TraitCategory string `db:"trait_category"`
}
