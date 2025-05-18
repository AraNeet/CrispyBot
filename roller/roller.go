package roller

import (
	"CrispyBot/database"
	"CrispyBot/variables"
	"math/rand"
	"time"
)

// Generate a new character with random attributes
func GenerateCharacter(ownerID string) database.Character {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate stats
	stats := generateStats(rng)

	// Generate attributes
	attributes := generateAttributes(rng)

	// Create the character
	character := database.Character{
		Owner:      ownerID,
		Stats:      stats,
		Attributes: attributes,
	}

	return character
}

// Generate random stats based on rarity
func generateStats(rng *rand.Rand) database.StatsSheets {

	// Generate each stat
	vitality := generateStat(variables.Vitality, VitalityRarity, rng)
	durability := generateStat(variables.Durability, DurabilityRarity, rng)
	speed := generateStat(variables.Speed, SpeedRarity, rng)
	strength := generateStat(variables.Strength, StrengthRarity, rng)
	intelligence := generateStat(variables.Intelligence, IntelligenceRarity, rng)
	manaFlow := generateStat(variables.ManaFlow, ManaFlowRarity, rng)
	skillLevel := generateStat(variables.SkillLevel, SkillLevelRarity, rng)

	return database.StatsSheets{
		Vitality:     vitality,
		Durability:   durability,
		Speed:        speed,
		Strength:     strength,
		Intelligence: intelligence,
		ManaFlow:     manaFlow,
		SkillLevel:   skillLevel,
	}
}

// Generate a single stat with random rarity
func generateStat(statType variables.StatType, rarityMap map[string][]string, rng *rand.Rand) database.Stat {
	rarity := SelectTier(config, rng)
	statName := RollRarityTrait(rarityMap, config, rng)

	// Get base value for the stat
	baseValue := getStatBaseValue(statType, statName)

	return database.Stat{
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
	case variables.ManaFlow:
		if val, ok := variables.ManaFlowValue[statName]; ok {
			value = val
		}
	case variables.SkillLevel:
		if val, ok := variables.SkillLevelValue[statName]; ok {
			value = val
		}
	}

	// Default to average value if not found
	if value == 0 {
		value = 127.5 // Default "Average" value
	}

	return int(value)
}

// Generate character attributes
func generateAttributes(rng *rand.Rand) database.Attributes {

	// Generate race trait
	raceTrait := generateRaceTrait(rng)

	// Generate element trait
	elementTrait := generateElementTrait(rng)

	// Generate extra trait
	extraTrait := generateExtraTrait(rng)

	// Generate weakness trait
	weaknessTrait := generateWeaknessTrait(rng)

	// Generate alignment trait
	alignmentTrait := generateAlignmentTrait(rng)

	// Generate x-factor trait
	xFactorTrait := generateXFactorTrait(rng)

	return database.Attributes{
		Race:      raceTrait,
		Element:   elementTrait,
		Trait:     extraTrait,
		Weakness:  weaknessTrait,
		Alignment: alignmentTrait,
		X_Factor:  xFactorTrait,
	}
}

// Generate a race trait
func generateRaceTrait(rng *rand.Rand) database.Trait {
	raceName := RollRarityTrait(RacesRarity, config, rng)
	rarity := getTierForTrait(raceName, RacesRarity)

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

	return database.Trait{
		Rarity:      rarity,
		Trait_Name:  raceName,
		Type:        variables.Race,
		Stats_Value: statsValues,
	}
}

// Generate an element trait
func generateElementTrait(rng *rand.Rand) database.Trait {
	elementName := RollWeightedOption(ElementOptions, rng)

	// Elements don't have stats values in the current schema,
	// but we'll prepare an empty map for future expansion
	statsValues := make(map[string]int)

	return database.Trait{
		Rarity:      "Common", // Default rarity for elements
		Trait_Name:  elementName,
		Type:        variables.Buff, // Elements are considered buffs
		Stats_Value: statsValues,
	}
}

// Generate an extra (buff) trait
func generateExtraTrait(rng *rand.Rand) database.Trait {
	traitName := RollRarityTrait(ExtraTraitRarity, config, rng)
	rarity := getTierForTrait(traitName, ExtraTraitRarity)

	// Get trait stat values
	statsValues := make(map[string]int)

	// Look up trait values
	if traitValuesArr, ok := variables.ExtraTraitValues[traitName]; ok && len(traitValuesArr) > 0 {
		for _, valueMap := range traitValuesArr {
			for statName, value := range valueMap {
				statsValues[statName] = int(value)
			}
		}
	}

	return database.Trait{
		Rarity:      rarity,
		Trait_Name:  traitName,
		Type:        variables.Buff,
		Stats_Value: statsValues,
	}
}

// Generate a weakness trait
func generateWeaknessTrait(rng *rand.Rand) database.Trait {
	weaknessName := RollWeightedOption(WeaknessOptions, rng)

	// Get weakness stat values
	statsValues := make(map[string]int)

	// Look up weakness values
	if weaknessValuesArr, ok := variables.WeaknessValues[weaknessName]; ok && len(weaknessValuesArr) > 0 {
		for _, valueMap := range weaknessValuesArr {
			for statName, value := range valueMap {
				// Weaknesses are negative stat modifiers
				statsValues[statName] = -int(value)
			}
		}
	}

	return database.Trait{
		Rarity:      "Common", // Default rarity for weaknesses
		Trait_Name:  weaknessName,
		Type:        variables.Weakness,
		Stats_Value: statsValues,
	}
}

// Generate an alignment trait
func generateAlignmentTrait(rng *rand.Rand) database.Trait {
	alignmentName := RollWeightedOption(AlignmentOptions, rng)

	// Alignments don't affect stats in the current schema
	statsValues := make(map[string]int)

	return database.Trait{
		Rarity:      "Common", // Default rarity for alignments
		Trait_Name:  alignmentName,
		Type:        variables.Buff, // Neutral trait type
		Stats_Value: statsValues,
	}
}

// Generate an x-factor trait
func generateXFactorTrait(rng *rand.Rand) database.Trait {
	xFactorName := RollWeightedOption(XFactorOptions, rng)

	// X-Factors don't have defined stat values in the schema yet
	statsValues := make(map[string]int)

	return database.Trait{
		Rarity:      "Rare", // Default rarity for x-factors
		Trait_Name:  xFactorName,
		Type:        variables.X_Factor,
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
