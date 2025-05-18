package roller

import (
	"math/rand"
)

var (
	config = DefaultRarityConfig()
)

func TierNames() []string {
	return []string{"Common", "Uncommon", "Rare", "Epic", "Legendary"}
}

func SelectTier(config RarityConfig, rng *rand.Rand) string {
	total := config.Common + config.Uncommon + config.Rare + config.Epic + config.Legendary
	roll := rng.Intn(total)
	if roll <= config.Common {
		return "Common"
	} else if roll <= config.Common+config.Uncommon {
		return "Uncommon"
	} else if roll <= config.Common+config.Uncommon+config.Rare {
		return "Rare"
	} else if roll <= config.Common+config.Uncommon+config.Rare+config.Epic {
		return "Epic"
	} else {
		return "Legendary"
	}
}

func RollRarityTrait(rarityMap map[string][]string, config RarityConfig, rng *rand.Rand) string {
	tier := SelectTier(config, rng)
	options, ok := rarityMap[tier]
	if !ok || len(options) == 0 {
		for _, t := range TierNames() {
			if opts, ok := rarityMap[t]; ok && len(opts) > 0 {
				return opts[rng.Intn(len(opts))]
			}
		}
		return ""
	}
	return options[rng.Intn(len(options))]
}

// Weighted roll function
func RollWeightedOption(options []WeightedOption, rng *rand.Rand) string {
	// 1. Calculate total weight
	totalWeight := 0
	for _, opt := range options {
		totalWeight += opt.Weight
	}

	if totalWeight <= 0 {
		return ""
	}
	roll := rng.Intn(totalWeight)

	// 3. Find the selected option
	for _, opt := range options {
		if roll < opt.Weight {
			return opt.Value
		}
		roll -= opt.Weight
	}

	// Fallback (should never hit if weights are correct)
	return options[0].Value
}

func RollEqualOption(options []string, rng *rand.Rand) string {
	if len(options) == 0 {
		return ""
	}
	idx := rng.Intn(len(options))
	return options[idx]
}
