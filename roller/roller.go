package roller

import (
	"CrispyBot/database/models"
	"CrispyBot/variables"
	"math/rand"
	"time"
)

// Generate a new character with random attributes
func GenerateCharacter(ownerID string) models.Character {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate stats
	stats := generateStats(rng)

	// Generate traits (innate, inadequacy, x-factor)
	traits := generateTraits(rng)

	// Generate characteristics (race, alignment, element)
	characteristics := generateCharacteristics(rng)

	// Create the character
	character := models.Character{
		Owner:           ownerID,
		Stats:           stats,
		Attributes:      traits,
		Characteristics: characteristics,
	}

	return character
}

// Generate random stats based on rarity
func generateStats(rng *rand.Rand) models.StatsSheets {
	// Generate each stat
	vitality := generateStat(variables.Vitality, VitalityRarity, rng)
	durability := generateStat(variables.Durability, DurabilityRarity, rng)
	speed := generateStat(variables.Speed, SpeedRarity, rng)
	strength := generateStat(variables.Strength, StrengthRarity, rng)
	intelligence := generateStat(variables.Intelligence, IntelligenceRarity, rng)
	mana := generateStat(variables.Mana, ManaRarity, rng)
	mastery := generateStat(variables.Mastery, MasteryRarity, rng)

	return models.StatsSheets{
		Vitality:     vitality,
		Durability:   durability,
		Speed:        speed,
		Strength:     strength,
		Intelligence: intelligence,
		Mana:         mana,
		Mastery:      mastery,
	}
}

// Generate a single stat with random rarity
func generateStat(statType variables.StatType, rarityMap map[string][]string, rng *rand.Rand) models.Stat {
	// First, select a trait name using RollRarityTrait
	statName := RollRarityTrait(rarityMap, config, rng)

	// Then determine the correct rarity based on the selected trait name
	rarity := getTierForTrait(statName, rarityMap)

	// Get base value for the stat
	baseValue := getStatBaseValue(statType, statName)

	return models.Stat{
		Rarity:    rarity,
		Stat_Name: statName,
		Type:      statType,
		Value:     baseValue,
	}
}

// Get the base value for a stat based on its name
func getStatBaseValue(statType variables.StatType, statName string) int {
	var value float64

	switch statType {
	case variables.Vitality:
		if val, ok := variables.VitalityValue[statName]; ok {
			value = val
		}
	case variables.Strength:
		if val, ok := variables.StrengthValue[statName]; ok {
			value = val
		}
	case variables.Speed:
		if val, ok := variables.SpeedValue[statName]; ok {
			value = val
		}
	case variables.Durability:
		if val, ok := variables.DurabilityValue[statName]; ok {
			value = val
		}
	case variables.Intelligence:
		if val, ok := variables.IntelligenceValue[statName]; ok {
			value = val
		}
	case variables.Mana:
		if val, ok := variables.ManaValue[statName]; ok {
			value = val
		}
	case variables.Mastery:
		if val, ok := variables.MasteryValue[statName]; ok {
			value = val
		}
	}

	// Default to average value if not found
	if value == 0 {
		value = 127.5 // Default "Average" value
	}

	return int(value)
}

// Generate character traits (innate, inadequacy, x-factor)
func generateTraits(rng *rand.Rand) models.Traits {
	// Generate innate trait (buff)
	innateTrait := generateInnateTrait(rng)

	// Generate inadequacy trait (weakness)
	inadequacyTrait := generateInadequacyTrait(rng)

	// Generate x-factor trait
	xFactorTrait := generateXFactorTrait(rng)

	return models.Traits{
		Innate:     innateTrait,
		Inadequacy: inadequacyTrait,
		X_Factor:   xFactorTrait,
	}
}

// Generate character characteristics (race, alignment, element)
func generateCharacteristics(rng *rand.Rand) models.Characteristics {
	// Generate race characteristic
	race := generateRaceCharacteristic(rng)

	// Generate alignment characteristic
	alignment := generateAlignmentCharacteristic(rng)

	// Generate element characteristic
	element := generateElementCharacteristic(rng)

	return models.Characteristics{
		Race:      race,
		Alignment: alignment,
		Element:   element,
	}
}

// Generate an innate trait
func generateInnateTrait(rng *rand.Rand) models.Trait {
	traitName := RollRarityTrait(InnateRarity, config, rng)
	rarity := getTierForTrait(traitName, InnateRarity)

	// Get trait stat values
	statsValues := make(map[string]int)

	// Look up trait values
	if traitValuesArr, ok := variables.InnateValues[traitName]; ok && len(traitValuesArr) > 0 {
		for _, valueMap := range traitValuesArr {
			for statName, value := range valueMap {
				statsValues[statName] = int(value)
			}
		}
	}

	return models.Trait{
		Rarity:      rarity,
		Trait_Name:  traitName,
		Type:        variables.Innate,
		Stats_Value: statsValues,
	}
}

// Generate an inadequacy trait
func generateInadequacyTrait(rng *rand.Rand) models.Trait {
	inadequacyName := RollWeightedOption(InadequacyOptions, rng)

	// Get trait stat values
	statsValues := make(map[string]int)

	// Look up inadequacy values
	if inadequacyValuesArr, ok := variables.InadequacyValues[inadequacyName]; ok && len(inadequacyValuesArr) > 0 {
		for _, valueMap := range inadequacyValuesArr {
			for statName, value := range valueMap {
				// Inadequacies are negative stat modifiers
				statsValues[statName] = -int(value)
			}
		}
	}

	return models.Trait{
		Rarity:      "Common", // Default rarity for inadequacies
		Trait_Name:  inadequacyName,
		Type:        variables.Inadequacy,
		Stats_Value: statsValues,
	}
}

// Generate an x-factor trait
func generateXFactorTrait(rng *rand.Rand) models.Trait {
	xFactorName := RollWeightedOption(XFactorOptions, rng)

	// X-Factors don't have defined stat values in the schema yet
	statsValues := make(map[string]int)

	return models.Trait{
		Rarity:      "Rare", // Default rarity for x-factors
		Trait_Name:  xFactorName,
		Type:        variables.X_Factor,
		Stats_Value: statsValues,
	}
}

// Generate a race characteristic
func generateRaceCharacteristic(rng *rand.Rand) models.Characteristic {
	raceName := RollRarityTrait(RaceRarity, config, rng)
	rarity := getTierForTrait(raceName, RaceRarity)

	// Get race stat values
	statsValues := make(map[string]int)

	// Look up race values in the variables
	if raceValuesArr, ok := variables.RaceValues[raceName]; ok && len(raceValuesArr) > 0 {
		// Race values are mapped as an array of maps with "Buff" and "Weakness" keys
		for _, raceValueMap := range raceValuesArr {
			// Process buffs
			if buffValues, ok := raceValueMap["Buff"]; ok && len(buffValues) > 0 {
				for _, buffMap := range buffValues {
					for statName, value := range buffMap {
						// Convert to int and add to stats
						statsValues[statName] = int(value)
					}
				}
			}

			// Process weaknesses (negative values)
			if weakValues, ok := raceValueMap["Weakness"]; ok && len(weakValues) > 0 {
				for _, weakMap := range weakValues {
					for statName, value := range weakMap {
						// Convert to int and subtract from stats
						statsValues[statName] = -int(value)
					}
				}
			}
		}
	}

	return models.Characteristic{
		Rarity:      rarity,
		Trait_Name:  raceName,
		Type:        variables.Race,
		Stats_Value: statsValues,
	}
}

// Generate an alignment characteristic
func generateAlignmentCharacteristic(rng *rand.Rand) models.Characteristic {
	alignmentName := RollWeightedOption(AlignmentOptions, rng)

	// Alignments don't affect stats in the current schema
	statsValues := make(map[string]int)

	return models.Characteristic{
		Rarity:      "Common", // Default rarity for alignments
		Trait_Name:  alignmentName,
		Type:        variables.Alignment,
		Stats_Value: statsValues,
	}
}

// Generate an element characteristic
func generateElementCharacteristic(rng *rand.Rand) models.Characteristic {
	elementName := RollWeightedOption(ElementOptions, rng)

	// Elements don't have stats values in the current schema,
	// but we'll prepare an empty map for future expansion
	statsValues := make(map[string]int)

	return models.Characteristic{
		Rarity:      "Common", // Default rarity for elements
		Trait_Name:  elementName,
		Type:        variables.Element,
		Stats_Value: statsValues,
	}
}

// Helper function to get the tier (rarity) for a trait
func getTierForTrait(traitName string, rarityMap map[string][]string) string {
	for tier, traits := range rarityMap {
		for _, trait := range traits {
			if trait == traitName {
				return tier
			}
		}
	}
	return "Common" // Default rarity
}
