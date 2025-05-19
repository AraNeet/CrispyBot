package combathandlers

import (
	"CrispyBot/variables"
	"fmt"
	"math/rand"
	"time"
)

// executeAction performs the selected action from the attacker to the target
func executeAction(attacker, target *CombatParticipant, actionName string) (string, error) {
	// Check if attacker is stunned
	if _, isStunned := attacker.StatusEffects["Stun"]; isStunned {
		return fmt.Sprintf("%s is stunned and cannot move!", attacker.UserName), nil
	}

	// Execute the appropriate action
	switch actionName {
	case "attack":
		return physicalAttack(attacker, target)
	case "magic":
		return magicalAttack(attacker, target)
	case "defend":
		return defend(attacker)
	case "item":
		return useItem(attacker)
	default:
		return "", fmt.Errorf("unknown action: %s", actionName)
	}
}

// physicalAttack executes a physical attack
func physicalAttack(attacker, target *CombatParticipant) (string, error) {
	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Check if attack hits
	hitChance := attacker.Accuracy
	hitRoll := rand.Intn(100)

	// Check for dodge
	dodgeChance := target.DodgeChance
	dodgeRoll := rand.Intn(100)

	// If dodge successful
	if dodgeRoll < dodgeChance {
		return fmt.Sprintf("%s attempts a physical attack, but %s dodges!", attacker.UserName, target.UserName), nil
	}

	// If attack misses
	if hitRoll >= hitChance {
		return fmt.Sprintf("%s's attack misses!", attacker.UserName), nil
	}

	// Calculate base damage
	damage := attacker.PhysicalDamage

	// Check for critical hit (base 5% chance)
	critChance := variables.BaseCritChance
	critRoll := rand.Intn(100)
	isCrit := critRoll < critChance

	if isCrit {
		damage = int(float64(damage) * variables.CritDamageMultiplier)
	}

	// Apply defense reduction
	defenseReduction := float64(target.Defense) / 100.0
	if defenseReduction > 0.75 {
		defenseReduction = 0.75 // Cap damage reduction at 75%
	}

	// Calculate final damage
	finalDamage := int(float64(damage) * (1.0 - defenseReduction))
	if finalDamage < 1 {
		finalDamage = 1 // Minimum damage is 1
	}

	// Apply damage
	target.CurrentHP -= finalDamage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}

	// Format result message
	var result string
	if isCrit {
		result = fmt.Sprintf("%s lands a critical hit for %d damage!", attacker.UserName, finalDamage)
	} else {
		result = fmt.Sprintf("%s attacks %s for %d damage!", attacker.UserName, target.UserName, finalDamage)
	}

	return result, nil
}

// magicalAttack executes a magical attack
func magicalAttack(attacker, target *CombatParticipant) (string, error) {
	// Check if attacker has enough mana
	manaCost := variables.MagicAttackBaseManaCost
	if attacker.CurrentMP < manaCost {
		return fmt.Sprintf("%s doesn't have enough mana to cast a spell!", attacker.UserName), nil
	}

	// Consume mana
	attacker.CurrentMP -= manaCost

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Check if spell hits
	hitChance := attacker.Accuracy - 5 // Magic is slightly harder to hit with
	hitRoll := rand.Intn(100)

	// Magic attacks can't be dodged as easily
	dodgeChance := target.DodgeChance / 2
	dodgeRoll := rand.Intn(100)

	// If dodge successful
	if dodgeRoll < dodgeChance {
		return fmt.Sprintf("%s casts a spell, but %s manages to avoid it!", attacker.UserName, target.UserName), nil
	}

	// If spell misses
	if hitRoll >= hitChance {
		return fmt.Sprintf("%s's spell fizzles out!", attacker.UserName), nil
	}

	// Calculate base damage
	damage := attacker.MagicalDamage

	// Check for critical hit (base 5% chance)
	critChance := variables.BaseCritChance
	critRoll := rand.Intn(100)
	isCrit := critRoll < critChance

	if isCrit {
		damage = int(float64(damage) * variables.CritDamageMultiplier)
	}

	// Apply elemental effectiveness
	effectiveness := getElementalEffectiveness(attacker.Element, target.Element)
	damage = int(float64(damage) * effectiveness)

	// Magic attacks ignore some defense
	defenseReduction := float64(target.Defense) / 200.0
	if defenseReduction > 0.5 {
		defenseReduction = 0.5 // Cap magical damage reduction at 50%
	}

	// Calculate final damage
	finalDamage := int(float64(damage) * (1.0 - defenseReduction))
	if finalDamage < 1 {
		finalDamage = 1 // Minimum damage is 1
	}

	// Apply damage
	target.CurrentHP -= finalDamage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}

	// Format result message with elemental effectiveness
	var result string
	if isCrit {
		result = fmt.Sprintf("%s casts a critical %s spell for %d damage!", attacker.UserName, attacker.Element, finalDamage)
	} else {
		result = fmt.Sprintf("%s casts a %s spell on %s for %d damage!", attacker.UserName, attacker.Element, target.UserName, finalDamage)
	}

	// Add effectiveness message
	if effectiveness > 1.0 {
		result += " It's super effective!"
	} else if effectiveness < 1.0 {
		result += " It's not very effective..."
	}

	// Chance to apply status effect based on element
	statusRoll := rand.Intn(100)
	if statusRoll < 20 { // 20% chance to apply status effect
		statusEffect := getElementalStatusEffect(attacker.Element)
		if statusEffect != "" {
			// Apply status effect (lasts 3 turns)
			target.StatusEffects[statusEffect] = 3
			result += fmt.Sprintf(" %s is now %s!", target.UserName, statusEffect)
		}
	}

	return result, nil
}

// defend increases defense for one turn
func defend(participant *CombatParticipant) (string, error) {
	// Increase defense by 50% for one turn
	defenseBoost := participant.Defense / 2
	participant.Defense += defenseBoost

	// Add status to remove boost next turn
	participant.StatusEffects["Defending"] = 1

	return fmt.Sprintf("%s takes a defensive stance, increasing defense!", participant.UserName), nil
}

// useItem uses an item from inventory (placeholder for now)
func useItem(participant *CombatParticipant) (string, error) {
	// Heal 20% of max HP
	healAmount := participant.MaxHP / 5
	participant.CurrentHP += healAmount

	// Cap at max HP
	if participant.CurrentHP > participant.MaxHP {
		participant.CurrentHP = participant.MaxHP
	}

	return fmt.Sprintf("%s uses a healing item, recovering %d HP!", participant.UserName, healAmount), nil
}

// getElementalEffectiveness returns the damage multiplier based on attacker and defender elements
func getElementalEffectiveness(attackerElement, defenderElement string) float64 {
	// Default to neutral effectiveness
	if attackerElement == "None" || defenderElement == "None" {
		return 1.0
	}

	// Check effectiveness table
	if effectivenessMap, exists := variables.ElementEffectiveness[attackerElement]; exists {
		if multiplier, exists := effectivenessMap[defenderElement]; exists {
			return multiplier
		}
	}

	// Default to neutral if not found in the table
	return 1.0
}

// getElementalStatusEffect returns a potential status effect based on element
func getElementalStatusEffect(element string) string {
	switch element {
	case "Fire":
		return "Burn"
	case "Toxic":
		return "Poison"
	case "Lightning":
		return "Stun"
	case "Frost":
		return "Freeze"
	default:
		return ""
	}
}
