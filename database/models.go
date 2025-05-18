package database

import (
	"CrispyBot/variables"
)

// Character Related Models

// User model
/*
	DiscordID - Main link to the user
	Character - User's Character
*/
type User struct {
	DiscordID string    `json:"DiscordID"`
	Character Character `json:"Character"`
}

// Character model
/*
	ID - Character ID
	Owner - Owner's ID
	Stats - Character main stats. EXM: Speed, Vitality, and Durability
	Attributes - Character main addition Attributes or traits. EXM: Race, Element, and Weakness
*/
type Character struct {
	ID         string      `json:"id,omitempty"` // Changed from owner to id
	Owner      string      `json:"owner"`        // Added explicit Owner field
	Stats      StatsSheets `json:"stats"`        // Ensure this matches DB field
	Attributes Attributes  `json:"attributes"`
}

// StatsSheets model
/*
	This holds all the main battle stats.
	ID - Stats ID which is the characters ID
	Vitality - Equals Health
	Durability - Equals Armor
	Speed - Equals Attack Order
	Strength - Equals Attack Power
	Intelligence - Equals Special Power
	ManaFlow - Equals Mana
	SkillLevel - Equals Attack/Special Efficiency
*/
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
/*
	This holds all the main attributes. These attributes affect main stats
	ID - Attributes ID which is the characters ID
	Race - Can be Negative or Positive Stats Boosts
	Element - Element affects character combat. For example if a Fire element user fight a Water user. The Water uses does more damage.const
	Trait - Traits are Positive Stats Boosts
	Weakness - Weakness Negative Stats Boosts
	Alignment - Alignment doesn't affect any stat
	X_Factor - X_Factors Boost Attributes or A stat. Its normally a attribute.
*/
type Attributes struct {
	Race      Trait `json:"race"`
	Element   Trait `json:"element"`
	Trait     Trait `json:"trait"`
	Weakness  Trait `json:"weakness"`
	Alignment Trait `json:"alignment"`
	X_Factor  Trait `json:"x_factor"`
}

// Traits model
/*
	This hold the values for traits. A trait is any of the attributes. but it with the stats it affects
	Rarity - Is how rare a trait is
	Trait_Name - The traits name
	Type - The trait type
	Stats_Value - It's the stats that are boosted or Reduced
*/
type Trait struct {
	Rarity      string              `json:"rarity"`
	Trait_Name  string              `json:"trait_name"`
	Type        variables.TraitType `json:"type"`
	Stats_Value map[string]int      `json:"stats_value"`
}

// Stats model
/*
	This hold the values for traits. A trait is any of the attributes. but it with the stats it affects
	Rarity - Is how rare a stat is
	Stat_Name - The stat name
	Type - The stat
	Value - stat value
*/
type Stat struct {
	Rarity    string             `json:"rarity"`
	Stat_Name string             `json:"stat_name"`
	Type      variables.StatType `json:"type"`
	Value     int                `json:"value"`
}
