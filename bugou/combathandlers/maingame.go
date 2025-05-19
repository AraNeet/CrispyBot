package combathandlers

import (
	"fmt"
	"sync"
	"time"
)

var (
	// ActiveBattleCleanup handles cleanup of old inactive battles
	ActiveBattleCleaner sync.Once
)

// InitializeCombatSystem sets up the combat system
func InitializeCombatSystem() {
	// Start background cleanup for stale battles
	ActiveBattleCleaner.Do(func() {
		go battleCleanupRoutine()
	})
}

// battleCleanupRoutine periodically checks for abandoned battles and removes them
func battleCleanupRoutine() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cleanupStaleBattles()
	}
}

// cleanupStaleBattles removes battles that have been inactive for too long
func cleanupStaleBattles() {
	now := time.Now()
	staleThreshold := 30 * time.Minute

	ActiveBattlesMutex.Lock()
	defer ActiveBattlesMutex.Unlock()

	staleBattleIDs := []string{}

	// Find stale battles
	for id, battle := range ActiveBattles {
		if now.Sub(battle.LastUpdated) > staleThreshold {
			staleBattleIDs = append(staleBattleIDs, id)
		}
	}

	// Remove stale battles
	for _, id := range staleBattleIDs {
		delete(ActiveBattles, id)
		fmt.Printf("Cleaned up stale battle: %s\n", id)
	}

	fmt.Printf("Battle cleanup: removed %d stale battles. Active battles: %d\n",
		len(staleBattleIDs), len(ActiveBattles))
}

// Available NPC types for battles
var NPCTemplates = map[string]int{
	"Training Dummy":   1,
	"Goblin":           2,
	"Bandit":           3,
	"Wolf Pack":        4,
	"Dark Knight":      5,
	"Troll":            6,
	"Dragon Whelp":     7,
	"Necromancer":      8,
	"Ancient Guardian": 9,
	"Dragon Lord":      10,
}

// GetNPCLevel returns the appropriate level for a named NPC type
func GetNPCLevel(npcType string, customDifficulty int) int {
	if level, exists := NPCTemplates[npcType]; exists {
		// If a custom difficulty was provided, use it
		if customDifficulty > 0 {
			return customDifficulty
		}
		return level
	}

	// Default to level 1 if NPC type not found
	return 1
}

// UpdateBattleTimestamp updates the last activity time for a battle
func UpdateBattleTimestamp(battleID string) {
	ActiveBattlesMutex.Lock()
	defer ActiveBattlesMutex.Unlock()

	if battle, exists := ActiveBattles[battleID]; exists {
		battle.LastUpdated = time.Now()
	}
}
