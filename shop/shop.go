package shop

import (
	"CrispyBot/database/models"
	"CrispyBot/roller"
	"math/rand"
	"time"
)

const (
	// Number of items to generate for the shop
	ShopInventorySize = 6

	// Rarity chances - similar to character generation
	CommonChance    = 50
	UncommonChance  = 25
	RareChance      = 15
	EpicChance      = 8
	LegendaryChance = 2

	// Base prices by rarity
	CommonPrice    = 100
	UncommonPrice  = 250
	RarePrice      = 500
	EpicPrice      = 1000
	LegendaryPrice = 2500

	// Stat buff/debuff values by rarity
	CommonStatValue    = 10
	UncommonStatValue  = 15
	RareStatValue      = 20
	EpicStatValue      = 30
	LegendaryStatValue = 60
)

// Available stats that can be buffed/debuffed
var statTypes = []string{
	"Vitality",
	"Durability",
	"Speed",
	"Strength",
	"Intelligence",
	"Mana",
	"Mastery",
}

// CreateShop generates a new shop with a random inventory
func CreateShop() models.Shop {
	// Set shop refresh timer to next midnight
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	// Create random inventory
	inventory := GenerateInventory()

	return models.Shop{
		Timer:     nextMidnight,
		Inventory: inventory,
	}
}

// IsShopExpired checks if the shop needs to be refreshed
func IsShopExpired(shop models.Shop) bool {
	return time.Now().After(shop.Timer)
}

// RefreshShop creates a new inventory and updates the timer
func RefreshShop(shop *models.Shop) {
	// Set new timer to next midnight
	now := time.Now()
	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	// Generate new inventory
	shop.Timer = nextMidnight
	shop.Inventory = GenerateInventory()
}

// generateInventory creates a random selection of items for the shop
func GenerateInventory() models.Inventory {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	items := make(map[int]models.Item)

	// Use the existing weapon options from the roller package
	weaponOptions := roller.WeaponOptions

	// Generate a random set of items
	for i := 0; i < ShopInventorySize; i++ {
		// Select a random weapon
		weaponName := roller.RollWeightedOption(weaponOptions, rng)

		// Generate random rarity for the item
		rarity := GenerateItemRarity(rng)

		// Generate random stats for the item based on rarity
		stats := GenerateItemStats(rarity, rng)

		// Calculate price based on rarity and stats
		price := CalculatePrice(rarity, stats)

		// Create and add the item to the inventory
		items[i] = models.Item{
			Name:   weaponName,
			Rarity: rarity,
			Stats:  stats,
			Price:  price,
		}
	}

	return models.Inventory{Items: items}
}

// generateItemRarity determines the rarity of an item
func GenerateItemRarity(rng *rand.Rand) string {
	rarityConfig := roller.RarityConfig{
		Common:    CommonChance,
		Uncommon:  UncommonChance,
		Rare:      RareChance,
		Epic:      EpicChance,
		Legendary: LegendaryChance,
	}

	return roller.SelectTier(rarityConfig, rng)
}

// generateItemStats creates random stat buffs/debuffs based on item rarity
func GenerateItemStats(rarity string, rng *rand.Rand) map[string]int {
	stats := make(map[string]int)

	// Determine base stat value based on rarity
	var baseValue int
	switch rarity {
	case "Common":
		baseValue = CommonStatValue
	case "Uncommon":
		baseValue = UncommonStatValue
	case "Rare":
		baseValue = RareStatValue
	case "Epic":
		baseValue = EpicStatValue
	case "Legendary":
		baseValue = LegendaryStatValue
	}

	// Determine how many stats will be affected (1-3 depending on rarity)
	numStats := 1
	if rarity == "Uncommon" || rarity == "Rare" {
		numStats = 1 + rng.Intn(2) // 1-2 stats
	} else if rarity == "Epic" || rarity == "Legendary" {
		numStats = 2 + rng.Intn(2) // 2-3 stats
	}

	// Select random stats and assign values
	availableStats := make([]string, len(statTypes))
	copy(availableStats, statTypes)

	for i := 0; i < numStats && len(availableStats) > 0; i++ {
		// Select a random stat
		statIdx := rng.Intn(len(availableStats))
		stat := availableStats[statIdx]

		// Remove the selected stat from available options
		availableStats = append(availableStats[:statIdx], availableStats[statIdx+1:]...)

		// 70% chance for positive buff, 30% chance for debuff
		value := baseValue
		if rng.Intn(10) < 3 {
			value = -value
		}

		stats[stat] = value
	}

	return stats
}

// calculatePrice determines an item's price based on rarity and stats
func CalculatePrice(rarity string, stats map[string]int) int {
	var basePrice int

	// Base price by rarity
	switch rarity {
	case "Common":
		basePrice = CommonPrice
	case "Uncommon":
		basePrice = UncommonPrice
	case "Rare":
		basePrice = RarePrice
	case "Epic":
		basePrice = EpicPrice
	case "Legendary":
		basePrice = LegendaryPrice
	}

	// Adjust price based on stats
	statModifier := 0
	for _, value := range stats {
		statModifier += value
	}

	// Positive stats increase price, negative stats decrease it
	priceModifier := 1.0 + float64(statModifier)/100.0
	if priceModifier < 0.5 {
		priceModifier = 0.5 // Minimum price is 50% of base price
	}

	return int(float64(basePrice) * priceModifier)
}
