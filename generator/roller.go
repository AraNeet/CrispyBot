package roller

import "CrispyBot/variables"

// Rarities Models

// Rarity Config
/*
	This is the base rarity config
*/
type RarityConfig struct {
	Common    int `json:"common"`
	Uncommon  int `json:"uncommon"`
	Rare      int `json:"rare"`
	Epic      int `json:"epic"`
	Legandary int `json:"legandary"`
}

func TierNames() []string {
	return []string{"Common", "Uncommon", "Rare", "Epic", "Legendary"}
}

func DefaultRarityConfig() RarityConfig {
	return RarityConfig{
		Common:    variables.Common_Chance,
		Uncommon:  variables.Uncommon_Chance,
		Rare:      variables.Rare_Chance,
		Epic:      variables.Epic_Chance,
		Legandary: variables.Legandary_Chance,
	}
}
