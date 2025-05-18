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
	Legendary int `json:"legendary"`
}

func DefaultRarityConfig() RarityConfig {
	return RarityConfig{
		Common:    variables.Common_Chance,
		Uncommon:  variables.Uncommon_Chance,
		Rare:      variables.Rare_Chance,
		Epic:      variables.Epic_Chance,
		Legendary: variables.Legendary_Chance,
	}
}
