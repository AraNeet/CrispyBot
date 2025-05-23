package models

import (
	"CrispyBot/variables"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Character Models
/*
	ID - ObjectID.
	Owner - Owners Discord ID.
	Chatacteristics - Character Apperance. Note: Can Boost or Nerf Stats.
	Stats - Character Stats.
	Traits - Positive/Negative Stats Boosts.
	EquippedWeapon - EquippedWeapon. Note: Can Boost or Nerf Stats.
	Level - Character level.
	Experience - How much until next level.
*/
type Character struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Owner           string             `bson:"Owner" json:"owner"`
	Characteristics Characteristics    `bson:"Characteriastics" json:"characteriastics"`
	Stats           StatsSheets        `bson:"Stats" json:"stats"`
	Traits          TraitsSheets       `bson:"Traits" json:"traits"`
	EquippedWeapon  EquippedItem       `bson:"EquippedWeapon" json:"equippedWeapon"`
	Level           int                `bson:"Level" json:"level"`
	Experience      int                `bson:"Experience" json:"experience"`
}

// Equipped Item Model
/*
	ItemKey - Unique identifier/key for the item type.
	ItemName - Display name of the equipped item.
*/
type EquippedItem struct {
	ItemKey  string `bson:"itemKey" json:"itemKey"`
	ItemName string `bson:"itemName" json:"itemName"`
}

// Item Record Model
/*
	ID - ObjectID for the item record.
	OwnerID - Discord ID of the item owner.
	InventoryKey - Unique key identifying this item in inventory.
	Item - The actual item data/properties.
	Timestamp - When this item was created/acquired.
*/
type ItemRecord struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID      string             `bson:"ownerID" json:"ownerID"`
	InventoryKey string             `bson:"inventoryKey" json:"inventoryKey"`
	Item         Item               `bson:"item" json:"item"`
	Timestamp    time.Time          `bson:"timestamp" json:"timestamp"`
}

// Stats Sheets Model
/*
	Vitality - Health/HP related stat.
	Durability - Defense/resistance stat.
	Speed - Movement/agility stat.
	Strength - Physical damage stat.
	Intelligence - Magical damage/wisdom stat.
	Mana - Magic points/energy stat.
	Mastery - Skill proficiency stat.
*/
type StatsSheets struct {
	Vitality     Stat `bson:"Vitality" json:"vitality"`
	Durability   Stat `bson:"Durability" json:"durability"`
	Speed        Stat `bson:"Speed" json:"speed"`
	Strength     Stat `bson:"Strength" json:"strength"`
	Intelligence Stat `bson:"Intelligence" json:"intelligence"`
	Mana         Stat `bson:"Mana" json:"mana"`
	Mastery      Stat `bson:"Mastery" json:"mastery"`
}

// Traits Sheets Model
/*
	Innate - Natural/born trait. Note: Usually positive modifiers.
	Inadequacy - Weakness/flaw trait. Note: Usually negative modifiers.
	X_Factor - Special/unique trait. Note: Can be positive or negative.
*/
type TraitsSheets struct {
	Innate     Trait `bson:"Innate" json:"innate"`
	Inadequacy Trait `bson:"Inadequacy" json:"inadequacy"`
	X_Factor   Trait `bson:"XFactor" json:"xfactor"`
}

// Characteristics Model
/*
	Race - Character's race/species. Note: Affects base stats.
	Alignment - Moral/ethical alignment. Note: May affect certain interactions.
	Element - Elemental affinity. Note: Affects damage types and resistances.
	Height - Physical height characteristic. Note: May affect certain stats.
*/
type Characteristics struct {
	Race      Characteristic `bson:"race" json:"race"`
	Alignment Characteristic `bson:"alignment" json:"alignment"`
	Element   Characteristic `bson:"element" json:"element"`
	Height    Characteristic `bson:"height" json:"height"`
}

// Individual Stat Model
/*
	Rarity - Rarity level of the stat (common, rare, epic, etc.).
	Stat_Name - Display name of the stat.
	Type - Stat type from variables enum.
	Value - Base stat value.
	EquipBonus - Bonus from equipped items.
	TraitBonus - Bonus from character traits.
	TotalValue - Final calculated stat value. Note: Sum of Value + EquipBonus + TraitBonus.
*/
type Stat struct {
	Rarity     string             `bson:"Rarity" json:"rarity"`
	Stat_Name  string             `bson:"StatName" json:"statName"`
	Type       variables.StatType `bson:"Type" json:"type"`
	Value      int                `bson:"Value" json:"value"`
	EquipBonus int                `bson:"EquipBonus" json:"equipBonus"`
	TraitBonus int                `json:"TraitBonus"`
	TotalValue int                `bson:"TotalValue" json:"totalValue"`
}

// Individual Trait Model
/*
	Rarity - Rarity level of the trait.
	Trait_Name - Display name of the trait.
	Type - Trait type from variables enum.
	Stats_Value - Map of stat modifications. Note: Key is stat name, value is modifier amount.
*/
type Trait struct {
	Rarity      string              `bson:"Rarity" json:"rarity"`
	Trait_Name  string              `bson:"TraitName" json:"traitName"`
	Type        variables.TraitType `bson:"Type" json:"type"`
	Stats_Value map[string]int      `bson:"StatsValue" json:"statsValue"`
}

// Individual Characteristic Model
/*
	Rarity - Rarity level of the characteristic.
	Trait_Name - Display name of the characteristic.
	Type - Characteristic type from variables enum.
	Stats_Value - Map of stat modifications. Note: Key is stat name, value is modifier amount.
*/
type Characteristic struct {
	Rarity      string                       `bson:"Rarity" json:"rarity"`
	Trait_Name  string                       `bson:"CharacteristicsName" json:"characteristicsName"`
	Type        variables.CharacteristicType `bson:"Type" json:"type"`
	Stats_Value map[string]int               `bson:"StatsValue" json:"statsValue"`
}
