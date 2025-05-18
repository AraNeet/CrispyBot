package roller

import (
	"math/rand"
	"testing"
	"time"
)

// Helper to create a new rng for each test
func newTestRNG() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestSelectTier_Distribution(t *testing.T) {
	config := RarityConfig{
		Common:    60,
		Uncommon:  20,
		Rare:      10,
		Epic:      7,
		Legendary: 3,
	}
	rng := newTestRNG()
	counts := map[string]int{}
	for i := 0; i < 10000; i++ {
		tier := SelectTier(config, rng)
		counts[tier]++
	}
	// Ensure all tiers are present at least once
	for _, tier := range TierNames() {
		if counts[tier] == 0 {
			t.Errorf("Tier %s was never selected", tier)
		}
	}
}

func TestRollRarityTrait_AllTiersPresent(t *testing.T) {
	config := RarityConfig{
		Common:    60,
		Uncommon:  20,
		Rare:      10,
		Epic:      7,
		Legendary: 3,
	}
	rng := newTestRNG()
	testMap := map[string][]string{
		"Common":    {"A"},
		"Uncommon":  {"B"},
		"Rare":      {"C"},
		"Epic":      {"D"},
		"Legendary": {"E"},
	}
	for i := 0; i < 100; i++ {
		trait := RollRarityTrait(testMap, config, rng)
		if trait == "" {
			t.Error("RollRarityTrait returned an empty string")
		}
	}
}

func TestRollRarityTrait_MissingTier(t *testing.T) {
	config := RarityConfig{
		Common:    60,
		Uncommon:  20,
		Rare:      10,
		Epic:      7,
		Legendary: 3,
	}
	rng := newTestRNG()
	// Only "Common" is present
	testMap := map[string][]string{
		"Common": {"A"},
	}
	for i := 0; i < 100; i++ {
		trait := RollRarityTrait(testMap, config, rng)
		if trait != "A" {
			t.Errorf("Expected 'A', got '%s'", trait)
		}
	}
}

func TestRollRarityTrait_EmptyMap(t *testing.T) {
	config := RarityConfig{
		Common:    60,
		Uncommon:  20,
		Rare:      10,
		Epic:      7,
		Legendary: 3,
	}
	rng := newTestRNG()
	testMap := map[string][]string{}
	trait := RollRarityTrait(testMap, config, rng)
	if trait != "" {
		t.Errorf("Expected empty string, got '%s'", trait)
	}
}

func TestRollWeightedOption_Normal(t *testing.T) {
	rng := newTestRNG()
	options := []WeightedOption{
		{Value: "A", Weight: 1},
		{Value: "B", Weight: 2},
		{Value: "C", Weight: 3},
	}
	for i := 0; i < 100; i++ {
		val := RollWeightedOption(options, rng)
		if val != "A" && val != "B" && val != "C" {
			t.Errorf("Unexpected value: %s", val)
		}
	}
}

func TestRollWeightedOption_Empty(t *testing.T) {
	rng := newTestRNG()
	options := []WeightedOption{}
	val := RollWeightedOption(options, rng)
	if val != "" && len(options) > 0 {
		t.Errorf("Expected empty string for empty options, got '%s'", val)
	}
}

func TestRollEqualOption_Normal(t *testing.T) {
	rng := newTestRNG()
	options := []string{"A", "B", "C"}
	for i := 0; i < 100; i++ {
		val := RollEqualOption(options, rng)
		if val != "A" && val != "B" && val != "C" {
			t.Errorf("Unexpected value: %s", val)
		}
	}
}

func TestRollEqualOption_Empty(t *testing.T) {
	rng := newTestRNG()
	options := []string{}
	val := RollEqualOption(options, rng)
	if val != "" {
		t.Errorf("Expected empty string for empty options, got '%s'", val)
	}
}
