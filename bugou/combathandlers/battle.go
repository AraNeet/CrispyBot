package combathandlers

import (
	"CrispyBot/database/models"
	"CrispyBot/variables"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Constants for battle state
const (
	BattlePending  = "pending"
	BattleOngoing  = "ongoing"
	BattleComplete = "complete"
)

// BattleResult represents the outcome of a battle
type BattleResult struct {
	Winner       string
	Loser        string
	Rounds       int
	WinnerHP     int
	Experience   int
	CurrencyGain int
}

// CombatParticipant represents a character with combat-ready stats
type CombatParticipant struct {
	Character      models.Character
	DiscordID      string
	UserName       string
	CurrentHP      int
	MaxHP          int
	CurrentMP      int
	MaxMP          int
	PhysicalDamage int
	MagicalDamage  int
	Defense        int
	Initiative     int
	Accuracy       int
	DodgeChance    int
	Element        string
	StatusEffects  map[string]int // Effect name -> remaining turns
	ActionThisTurn string
	TargetThisTurn string
	IsBot          bool // Flag for NPC opponents
}

// Battle represents a combat encounter between two participants
type Battle struct {
	ID                 string
	ChannelID          string
	Participants       map[string]*CombatParticipant
	CurrentTurn        string
	Round              int
	State              string
	LastUpdated        time.Time
	TurnOrder          []string // IDs in initiative order
	Log                []string // Combat log
	InteractionMessage string   // Discord message ID for battle UI
}

// NewBattle initializes a new battle between two participants
func NewBattle(channelID string, participant1 *CombatParticipant, participant2 *CombatParticipant) *Battle {
	battleID := fmt.Sprintf("battle_%s_%s_%d", participant1.DiscordID, participant2.DiscordID, time.Now().Unix())

	// Initialize participants map
	participants := make(map[string]*CombatParticipant)
	participants[participant1.DiscordID] = participant1
	participants[participant2.DiscordID] = participant2

	// Determine turn order based on initiative (speed)
	turnOrder := determineTurnOrder(participant1, participant2)

	return &Battle{
		ID:           battleID,
		ChannelID:    channelID,
		Participants: participants,
		CurrentTurn:  turnOrder[0], // First in initiative order goes first
		Round:        1,
		State:        BattlePending,
		LastUpdated:  time.Now(),
		TurnOrder:    turnOrder,
		Log:          []string{fmt.Sprintf("Battle between %s and %s begins!", participant1.UserName, participant2.UserName)},
	}
}

// CharacterToCombatParticipant converts a character to a combat-ready participant
func CharacterToCombatParticipant(character models.Character, discordID string, userName string) *CombatParticipant {
	// Apply stat caps
	cappedStats := capStats(character.Stats)

	// Calculate combat stats with 10:1 ratio
	maxHP := cappedStats.Vitality.TotalValue * variables.VitalityToHPRatio
	maxMP := cappedStats.Mana.TotalValue * variables.ManaToPoolRatio
	physDamage := cappedStats.Strength.TotalValue * variables.StrengthToDamageRatio
	magDamage := cappedStats.Intelligence.TotalValue * variables.IntelligenceToDamageRatio
	defense := cappedStats.Durability.TotalValue * variables.DurabilityToDefenseRatio
	initiative := cappedStats.Speed.TotalValue * variables.SpeedToInitiativeRatio

	// Calculate accuracy (base + mastery bonus)
	accuracy := variables.BaseAccuracy + (cappedStats.Mastery.TotalValue * variables.MasteryToAccuracyRatio / 10)

	// Calculate dodge chance (base + speed bonus, capped at 30%)
	dodgeChance := variables.BaseDodgeChance + (cappedStats.Speed.TotalValue / 10)
	if dodgeChance > variables.MaxDodgeChance {
		dodgeChance = variables.MaxDodgeChance
	}

	// Get element from character
	element := character.Characteristics.Element.Trait_Name

	return &CombatParticipant{
		Character:      character,
		DiscordID:      discordID,
		UserName:       userName,
		CurrentHP:      maxHP,
		MaxHP:          maxHP,
		CurrentMP:      maxMP,
		MaxMP:          maxMP,
		PhysicalDamage: physDamage,
		MagicalDamage:  magDamage,
		Defense:        defense,
		Initiative:     initiative,
		Accuracy:       accuracy,
		DodgeChance:    dodgeChance,
		Element:        element,
		StatusEffects:  make(map[string]int),
		IsBot:          false,
	}
}

// CreateNPCOpponent creates a computer-controlled opponent with the given stats
func CreateNPCOpponent(name string, level int) *CombatParticipant {
	// Scale stats based on level
	baseValue := 50 + (level * 5)
	if baseValue > variables.MaxStatValue {
		baseValue = variables.MaxStatValue
	}

	// Create NPC stats
	maxHP := baseValue * variables.VitalityToHPRatio
	maxMP := baseValue * variables.ManaToPoolRatio
	physDamage := baseValue * variables.StrengthToDamageRatio
	magDamage := baseValue * variables.IntelligenceToDamageRatio
	defense := baseValue * variables.DurabilityToDefenseRatio
	initiative := baseValue * variables.SpeedToInitiativeRatio

	// Calculate accuracy and dodge
	accuracy := variables.BaseAccuracy + (baseValue / 3)
	dodgeChance := variables.BaseDodgeChance + (baseValue / 10)
	if dodgeChance > variables.MaxDodgeChance {
		dodgeChance = variables.MaxDodgeChance
	}

	// Randomly select element
	elements := []string{"Fire", "Water", "Earth", "Wind", "Nature", "Lightning", "Ice", "Dark", "Light"}
	rand.Seed(time.Now().UnixNano())
	element := elements[rand.Intn(len(elements))]

	return &CombatParticipant{
		DiscordID:      "npc_" + fmt.Sprintf("%d", time.Now().UnixNano()),
		UserName:       name,
		CurrentHP:      maxHP,
		MaxHP:          maxHP,
		CurrentMP:      maxMP,
		MaxMP:          maxMP,
		PhysicalDamage: physDamage,
		MagicalDamage:  magDamage,
		Defense:        defense,
		Initiative:     initiative,
		Accuracy:       accuracy,
		DodgeChance:    dodgeChance,
		Element:        element,
		StatusEffects:  make(map[string]int),
		IsBot:          true,
	}
}

// capStats ensures no stat exceeds the maximum allowed value
func capStats(stats models.StatsSheets) models.StatsSheets {
	// Helper function to cap a single stat
	capStat := func(value int) int {
		if value > variables.MaxStatValue {
			return variables.MaxStatValue
		}
		return value
	}

	// Cap each stat's TotalValue
	stats.Vitality.TotalValue = capStat(stats.Vitality.TotalValue)
	stats.Durability.TotalValue = capStat(stats.Durability.TotalValue)
	stats.Speed.TotalValue = capStat(stats.Speed.TotalValue)
	stats.Strength.TotalValue = capStat(stats.Strength.TotalValue)
	stats.Intelligence.TotalValue = capStat(stats.Intelligence.TotalValue)
	stats.Mana.TotalValue = capStat(stats.Mana.TotalValue)
	stats.Mastery.TotalValue = capStat(stats.Mastery.TotalValue)

	return stats
}

// determineTurnOrder sets the initiative order based on speed
func determineTurnOrder(p1, p2 *CombatParticipant) []string {
	if p1.Initiative > p2.Initiative {
		return []string{p1.DiscordID, p2.DiscordID}
	} else if p2.Initiative > p1.Initiative {
		return []string{p2.DiscordID, p1.DiscordID}
	} else {
		// In case of tie, randomly determine who goes first
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(2) == 0 {
			return []string{p1.DiscordID, p2.DiscordID}
		} else {
			return []string{p2.DiscordID, p1.DiscordID}
		}
	}
}

// ProcessTurn executes the current participant's action
func (b *Battle) ProcessTurn() (string, error) {
	if b.State != BattleOngoing {
		return "", errors.New("battle is not in progress")
	}

	// Get the current participant
	currentParticipant := b.Participants[b.CurrentTurn]

	// Skip turn if no action set (shouldn't happen normally)
	if currentParticipant.ActionThisTurn == "" {
		return "No action selected", nil
	}

	// Get target
	targetID := currentParticipant.TargetThisTurn
	target, exists := b.Participants[targetID]
	if !exists {
		return "", errors.New("target not found")
	}

	// Process status effects at start of turn
	processStatusEffects(currentParticipant)

	// Execute the selected action
	result, err := executeAction(currentParticipant, target, currentParticipant.ActionThisTurn)
	if err != nil {
		return "", err
	}

	// Log the result
	b.Log = append(b.Log, result)

	// Check if battle is over
	if target.CurrentHP <= 0 {
		target.CurrentHP = 0
		b.State = BattleComplete
		b.Log = append(b.Log, fmt.Sprintf("%s has been defeated! %s wins the battle!", target.UserName, currentParticipant.UserName))
		return result, nil
	}

	// Move to next turn
	b.advanceTurn()

	return result, nil
}

// advanceTurn moves to the next participant in turn order
func (b *Battle) advanceTurn() {
	// Find current position in turn order
	var currentPos int
	for i, id := range b.TurnOrder {
		if id == b.CurrentTurn {
			currentPos = i
			break
		}
	}

	// Move to next participant
	nextPos := (currentPos + 1) % len(b.TurnOrder)

	// If we've gone through all participants, increment round counter
	if nextPos == 0 {
		b.Round++

		// Process end-of-round effects
		for _, participant := range b.Participants {
			// Reduce status effect durations
			for effect, turns := range participant.StatusEffects {
				if turns > 0 {
					participant.StatusEffects[effect] = turns - 1
				}
				if participant.StatusEffects[effect] == 0 {
					delete(participant.StatusEffects, effect)
				}
			}
		}
	}

	b.CurrentTurn = b.TurnOrder[nextPos]

	// Reset action selection for next turn
	nextParticipant := b.Participants[b.CurrentTurn]
	nextParticipant.ActionThisTurn = ""
	nextParticipant.TargetThisTurn = ""

	// If next is bot/NPC, auto-select its action
	if nextParticipant.IsBot {
		selectNPCAction(b, nextParticipant)
	}
}

// selectNPCAction chooses an action for an NPC
func selectNPCAction(battle *Battle, npc *CombatParticipant) {
	// Simple AI: choose between physical and magical attack based on stats
	var action string

	// Find target (the human player)
	var targetID string
	for id := range battle.Participants {
		if id != npc.DiscordID {
			targetID = id
			break
		}
	}

	// Choose action based on stronger stat
	if npc.PhysicalDamage > npc.MagicalDamage {
		action = "attack"
	} else if npc.CurrentMP >= variables.MagicAttackBaseManaCost {
		action = "magic"
	} else {
		action = "attack"
	}

	// Set NPC's action and target
	npc.ActionThisTurn = action
	npc.TargetThisTurn = targetID
}

// processStatusEffects applies effects of status conditions
func processStatusEffects(participant *CombatParticipant) {
	for effect, _ := range participant.StatusEffects {
		switch effect {
		case "Burn":
			// Burn does damage equal to 5% of max HP
			damage := participant.MaxHP / 20
			participant.CurrentHP -= damage
			if participant.CurrentHP < 0 {
				participant.CurrentHP = 0
			}
		case "Poison":
			// Poison does increasing damage each turn
			damage := participant.MaxHP / 10
			participant.CurrentHP -= damage
			if participant.CurrentHP < 0 {
				participant.CurrentHP = 0
			}
		case "Stun":
			// Stunned participants skip their turn (handled in executeAction)
		}
	}
}

// GetBattleStatus returns a formatted status of the current battle
func (b *Battle) GetBattleStatus() string {
	var status string

	// Get both participants
	var p1, p2 *CombatParticipant
	for _, p := range b.Participants {
		if p1 == nil {
			p1 = p
		} else {
			p2 = p
		}
	}

	status += fmt.Sprintf("⚔️ Round %d ⚔️\n\n", b.Round)

	// Show participant health and mana
	status += fmt.Sprintf("%s: HP %d/%d | MP %d/%d", p1.UserName, p1.CurrentHP, p1.MaxHP, p1.CurrentMP, p1.MaxMP)
	if len(p1.StatusEffects) > 0 {
		status += " | Status: "
		for effect, turns := range p1.StatusEffects {
			status += fmt.Sprintf("%s (%d) ", effect, turns)
		}
	}
	status += "\n"

	status += fmt.Sprintf("%s: HP %d/%d | MP %d/%d", p2.UserName, p2.CurrentHP, p2.MaxHP, p2.CurrentMP, p2.MaxMP)
	if len(p2.StatusEffects) > 0 {
		status += " | Status: "
		for effect, turns := range p2.StatusEffects {
			status += fmt.Sprintf("%s (%d) ", effect, turns)
		}
	}
	status += "\n\n"

	// Show whose turn it is
	currentParticipant := b.Participants[b.CurrentTurn]
	status += fmt.Sprintf("Current turn: %s\n", currentParticipant.UserName)

	// Show recent battle log (last 3 entries)
	status += "\nRecent actions:\n"
	startIdx := len(b.Log) - 3
	if startIdx < 0 {
		startIdx = 0
	}
	for i := startIdx; i < len(b.Log); i++ {
		status += "• " + b.Log[i] + "\n"
	}

	return status
}

// SetAction sets a participant's action for their turn
func (b *Battle) SetAction(userID string, action string, targetID string) error {
	// Verify it's this user's turn
	if b.CurrentTurn != userID {
		return errors.New("it's not your turn")
	}

	// Verify action is valid
	validActions := []string{"attack", "magic", "defend", "item"}
	actionValid := false
	for _, a := range validActions {
		if action == a {
			actionValid = true
			break
		}
	}

	if !actionValid {
		return errors.New("invalid action")
	}

	// Verify target is valid
	if _, exists := b.Participants[targetID]; !exists {
		return errors.New("invalid target")
	}

	// Set the action
	participant := b.Participants[userID]
	participant.ActionThisTurn = action
	participant.TargetThisTurn = targetID

	return nil
}

// StartBattle begins the battle
func (b *Battle) StartBattle() {
	b.State = BattleOngoing

	// If first participant is a bot, choose its action
	firstParticipant := b.Participants[b.CurrentTurn]
	if firstParticipant.IsBot {
		selectNPCAction(b, firstParticipant)
	}
}

func (b *Battle) GetResult() (*BattleResult, error) {
	if b.State != BattleComplete {
		return nil, errors.New("battle is not complete")
	}

	// Determine winner and loser
	var winnerID, loserID string
	for id, participant := range b.Participants {
		if participant.CurrentHP > 0 {
			winnerID = id
		} else {
			loserID = id
		}
	}

	// Calculate rewards based on battle stats
	winner := b.Participants[winnerID]

	// Base XP scaled by round count and modified by XP modifier
	baseExpGain := variables.BaseExperienceGain + (b.Round * 10)
	expGain := int(float64(baseExpGain) * variables.ExperienceModifier)

	// Base currency reward
	currencyGain := 100 + (b.Round * 5)

	return &BattleResult{
		Winner:       winnerID,
		Loser:        loserID,
		Rounds:       b.Round,
		WinnerHP:     winner.CurrentHP,
		Experience:   expGain,
		CurrencyGain: currencyGain,
	}, nil
}
