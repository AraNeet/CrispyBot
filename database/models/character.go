package models

import (
	"CrispyBot/variables"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Character struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner           string             `bson:"owner" json:"owner"`
	Characteristics Characteristics    `bson:"characteriastics" json:"characteriastics"`
	Stats           StatsSheets        `bson:"stats" json:"stats"`
	Attributes      Traits             `bson:"attributes" json:"attributes"`
}

type StatsSheets struct {
	Vitality     Stat `bson:"vitality" json:"vitality"`
	Durability   Stat `bson:"durability" json:"durability"`
	Speed        Stat `bson:"speed" json:"speed"`
	Strength     Stat `bson:"strength" json:"strength"`
	Intelligence Stat `bson:"intelligence" json:"intelligence"`
	Mana         Stat `bson:"manaflow" json:"manaflow"`
	Mastery      Stat `bson:"skillLevel" json:"skillLevel"`
}

type Traits struct {
	Innate     Trait `bson:"trait" json:"trait"`
	Inadequacy Trait `bson:"weakness" json:"weakness"`
	X_Factor   Trait `bson:"xFactor" json:"xFactor"`
}

type Characteristics struct {
	Race      Characteristic `bson:"race" json:"race"`
	Alignment Characteristic `bson:"alignment" json:"alignment"`
	Element   Characteristic `bson:"element" json:"element"`
}

type Stat struct {
	Rarity    string             `bson:"rarity" json:"rarity"`
	Stat_Name string             `bson:"statName" json:"statName"`
	Type      variables.StatType `bson:"type" json:"type"`
	Value     int                `bson:"value" json:"value"`
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
