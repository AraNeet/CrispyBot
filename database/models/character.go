package models

import (
	"CrispyBot/variables"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// In database/models/character.go
type Character struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner           string             `bson:"owner" json:"owner"`
	Characteristics Characteristics    `bson:"characteriastics" json:"characteriastics"`
	Stats           StatsSheets        `bson:"stats" json:"stats"`
	Attributes      Traits             `bson:"attributes" json:"attributes"`
	EquippedWeapon  EquippedItem       `bson:"equippedWeapon" json:"equippedWeapon"`
	Level           int                `bson:"level" json:"level"`           // New field for level
	Experience      int                `bson:"experience" json:"experience"` // New field for XP
}

type EquippedItem struct {
	ItemKey  string `bson:"itemKey" json:"itemKey"`
	ItemName string `bson:"itemName" json:"itemName"`
}

type ItemRecord struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID      string             `bson:"ownerID" json:"ownerID"`
	InventoryKey string             `bson:"inventoryKey" json:"inventoryKey"`
	Item         Item               `bson:"item" json:"item"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
}

type StatsSheets struct {
	Vitality     Stat `bson:"vitality" json:"vitality"`
	Durability   Stat `bson:"durability" json:"durability"`
	Speed        Stat `bson:"speed" json:"speed"`
	Strength     Stat `bson:"strength" json:"strength"`
	Intelligence Stat `bson:"intelligence" json:"intelligence"`
	Mana         Stat `bson:"Mana" json:"Mana"`
	Mastery      Stat `bson:"mastery" json:"mastery"`
}

type Traits struct {
	Innate     Trait `bson:"trait" json:"trait"`
	Inadequacy Trait `bson:"weakness" json:"weakness"`
	X_Factor   Trait `bson:"xFactor" json:"xFactor"`
}

// Update Characteristics to explicitly include Height
type Characteristics struct {
	Race      Characteristic `bson:"race" json:"race"`
	Alignment Characteristic `bson:"alignment" json:"alignment"`
	Element   Characteristic `bson:"element" json:"element"`
	Height    Characteristic `bson:"height" json:"height"` // Add Height field
}

type Stat struct {
	Rarity     string             `bson:"rarity" json:"rarity"`
	Stat_Name  string             `bson:"statName" json:"statName"`
	Type       variables.StatType `bson:"type" json:"type"`
	Value      int                `bson:"value" json:"value"`
	EquipBonus int                `bson:"equipBonus" json:"equipBonus"`
	TraitBonus int                `json:"traitBonus"` // New field for trait bonuses (not stored in DB)
	TotalValue int                `bson:"totalValue" json:"totalValue"`
}

type Trait struct {
	Rarity      string              `bson:"rarity" json:"rarity"`
	Trait_Name  string              `bson:"traitName" json:"traitName"`
	Type        variables.TraitType `bson:"type" json:"type"`
	Stats_Value map[string]int      `bson:"statsValue" json:"statsValue"`
}

type Characteristic struct {
	Rarity      string                       `bson:"rarity" json:"rarity"`
	Trait_Name  string                       `bson:"CharacteristicsName" json:"CharacteristicsName"`
	Type        variables.CharacteristicType `bson:"type" json:"type"`
	Stats_Value map[string]int               `bson:"statsValue" json:"statsValue"`
}
