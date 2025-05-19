package combathandlers

import (
	"CrispyBot/database"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// Map to track active battles
	ActiveBattles      = make(map[string]*Battle)
	ActiveBattlesMutex sync.Mutex
)

// HandleBattleCommand processes battle-related commands
func HandleBattleCommand(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		session.ChannelMessageSend(message.ChannelID, "Invalid battle command. Usage: `!cb battle [start|attack|magic|defend|item]`")
		return
	}

	subCommand := strings.ToLower(args[2])

	switch subCommand {
	case "start":
		// Start a battle with NPC or another player
		if len(args) >= 4 {
			// Check if it's a PvP request
			if strings.HasPrefix(args[3], "<@") && strings.HasSuffix(args[3], ">") {
				// Extract target user ID
				targetID := strings.TrimPrefix(strings.TrimSuffix(args[3], ">"), "<@")
				handlePvPBattleRequest(session, message, targetID)
			} else {
				// Start battle with NPC
				difficulty := 1
				if len(args) >= 5 {
					// Try to parse difficulty level
					_, err := fmt.Sscanf(args[4], "%d", &difficulty)
					if err != nil || difficulty < 1 || difficulty > 10 {
						difficulty = 1
					}
				}
				handleNPCBattle(session, message, args[3], difficulty)
			}
		} else {
			// Default: start a battle with an NPC
			handleNPCBattle(session, message, "Training Dummy", 1)
		}

	case "attack", "magic", "defend", "item":
		// Execute combat action
		handleCombatAction(session, message, subCommand)

	case "status":
		// Show battle status
		showBattleStatus(session, message)

	case "forfeit":
		// Forfeit battle
		forfeitBattle(session, message)

	default:
		session.ChannelMessageSend(message.ChannelID, "Unknown battle command. Available commands: start, attack, magic, defend, item, status, forfeit")
	}
}

// handleNPCBattle starts a battle with an NPC
func handleNPCBattle(session *discordgo.Session, message *discordgo.MessageCreate, npcName string, difficulty int) {
	// Get the database singleton
	db := database.GetDB()

	// Get user's character
	character, err := database.GetCharacterByOwner(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("You need a character to battle! Use `!cb roll` to create one."))
		return
	}

	// Check if user is already in a battle
	ActiveBattlesMutex.Lock()
	for _, battle := range ActiveBattles {
		if _, exists := battle.Participants[message.Author.ID]; exists {
			ActiveBattlesMutex.Unlock()
			session.ChannelMessageSend(message.ChannelID, "You are already in a battle! Finish or forfeit it first.")
			return
		}
	}
	ActiveBattlesMutex.Unlock()

	// Create combat participant from character
	player := CharacterToCombatParticipant(character, message.Author.ID, message.Author.Username)

	// Create NPC opponent
	npc := CreateNPCOpponent(npcName, difficulty)

	// Create the battle
	battle := NewBattle(message.ChannelID, player, npc)

	// Store battle in active battles map
	ActiveBattlesMutex.Lock()
	ActiveBattles[battle.ID] = battle
	ActiveBattlesMutex.Unlock()

	// Start the battle
	battle.StartBattle()

	// Create the battle embed
	battleEmbed := createBattleEmbed(battle)

	// Send battle start message
	msg, err := session.ChannelMessageSendEmbed(message.ChannelID, battleEmbed)
	if err != nil {
		fmt.Printf("Error sending battle message: %v\n", err)
		return
	}

	// Store message ID for updates
	battle.InteractionMessage = msg.ID

	// If NPC goes first, process their turn automatically
	if battle.CurrentTurn == npc.DiscordID {
		processBotTurn(session, battle)
	}
}

// handlePvPBattleRequest sends a battle challenge to another player
func handlePvPBattleRequest(session *discordgo.Session, message *discordgo.MessageCreate, targetID string) {
	// Get the database singleton
	db := database.GetDB()

	// Check if target is valid
	_, err := database.GetCharacterByOwner(db, targetID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "That user doesn't have a character!")
		return
	}

	// Check if either player is already in a battle
	ActiveBattlesMutex.Lock()
	defer ActiveBattlesMutex.Unlock()

	for _, battle := range ActiveBattles {
		if _, exists := battle.Participants[message.Author.ID]; exists {
			session.ChannelMessageSend(message.ChannelID, "You are already in a battle! Finish or forfeit it first.")
			return
		}
		if _, exists := battle.Participants[targetID]; exists {
			session.ChannelMessageSend(message.ChannelID, "That player is already in a battle!")
			return
		}
	}

	// Create PvP battle request embed
	challengeEmbed := &discordgo.MessageEmbed{
		Title:       "⚔️ Battle Challenge!",
		Description: fmt.Sprintf("<@%s> has challenged <@%s> to a battle!", message.Author.ID, targetID),
		Color:       0xFF0000,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "How to Accept",
				Value: "Click the Accept button below to begin the battle.",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Challenge expires in 60 seconds",
		},
	}

	// Add buttons for accepting/declining
	acceptButton := discordgo.Button{
		Label:    "Accept Challenge",
		Style:    discordgo.SuccessButton,
		CustomID: fmt.Sprintf("battle_accept_%s_%s", message.Author.ID, targetID),
	}

	declineButton := discordgo.Button{
		Label:    "Decline",
		Style:    discordgo.DangerButton,
		CustomID: fmt.Sprintf("battle_decline_%s_%s", message.Author.ID, targetID),
	}

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{acceptButton, declineButton},
	}

	// Send challenge message
	_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Embed:      challengeEmbed,
		Components: []discordgo.MessageComponent{actionRow},
	})

	if err != nil {
		fmt.Printf("Error sending PvP challenge: %v\n", err)
		return
	}

	// Set up handlers for button interactions and timeout
	// Note: This is simplified for brevity, in a real implementation you'd
	// need more robust interaction handling

	// Set up interaction handler for the buttons
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		// Handle battle acceptance
		if i.MessageComponentData().CustomID == fmt.Sprintf("battle_accept_%s_%s", message.Author.ID, targetID) {
			if i.Member.User.ID == targetID {
				// Start the PvP battle
				startPvPBattle(s, message.Author.ID, targetID, message.ChannelID, i.Message.ID)

				// Respond to the interaction
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Content:    "Battle challenge accepted! Starting battle...",
						Components: []discordgo.MessageComponent{},
					},
				})
			} else {
				// Only the challenged player can accept
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This challenge isn't for you!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		} else if i.MessageComponentData().CustomID == fmt.Sprintf("battle_decline_%s_%s", message.Author.ID, targetID) {
			if i.Member.User.ID == targetID {
				// Challenged player declined
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Content:    fmt.Sprintf("<@%s> declined the battle challenge.", targetID),
						Components: []discordgo.MessageComponent{},
					},
				})
			} else {
				// Only the challenged player can decline
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "This challenge isn't for you!",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		}
	})

	// Add timeout to remove buttons after 60 seconds
	time.AfterFunc(60*time.Second, func() {
		session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    message.ChannelID,
			ID:         message.ID,
			Content:    "Battle challenge expired.",
			Components: []discordgo.MessageComponent{},
		})
	})
}

// startPvPBattle creates a battle between two players
func startPvPBattle(session *discordgo.Session, player1ID, player2ID, channelID, messageID string) {
	// Get the database singleton
	db := database.GetDB()

	// Get both characters
	char1, err1 := database.GetCharacterByOwner(db, player1ID)
	char2, err2 := database.GetCharacterByOwner(db, player2ID)

	if err1 != nil || err2 != nil {
		session.ChannelMessageSend(channelID, "Error starting battle: One or both players don't have characters.")
		return
	}

	// Get usernames
	member1, err := session.GuildMember(channelID, player1ID)
	member2, err := session.GuildMember(channelID, player2ID)

	username1 := player1ID
	username2 := player2ID

	if err == nil && member1 != nil {
		if member1.Nick != "" {
			username1 = member1.Nick
		} else {
			username1 = member1.User.Username
		}
	}

	if err == nil && member2 != nil {
		if member2.Nick != "" {
			username2 = member2.Nick
		} else {
			username2 = member2.User.Username
		}
	}

	// Create combat participants
	p1 := CharacterToCombatParticipant(char1, player1ID, username1)
	p2 := CharacterToCombatParticipant(char2, player2ID, username2)

	// Create the battle
	battle := NewBattle(channelID, p1, p2)

	// Store battle in active battles map
	ActiveBattlesMutex.Lock()
	ActiveBattles[battle.ID] = battle
	ActiveBattlesMutex.Unlock()

	// Start the battle
	battle.StartBattle()

	// Create the battle embed
	battleEmbed := createBattleEmbed(battle)

	// Send battle start message
	msg, err := session.ChannelMessageSendEmbed(channelID, battleEmbed)
	if err != nil {
		fmt.Printf("Error sending battle message: %v\n", err)
		return
	}

	// Store message ID for updates
	battle.InteractionMessage = msg.ID
}

// handleCombatAction processes a player's combat action
func handleCombatAction(session *discordgo.Session, message *discordgo.MessageCreate, actionName string) {
	// Find the battle this player is in
	var playerBattle *Battle
	var battleID string

	ActiveBattlesMutex.Lock()
	for id, battle := range ActiveBattles {
		if _, exists := battle.Participants[message.Author.ID]; exists {
			playerBattle = battle
			battleID = id
			break
		}
	}
	ActiveBattlesMutex.Unlock()

	if playerBattle == nil {
		session.ChannelMessageSend(message.ChannelID, "You're not in a battle! Use `!cb battle start` to begin one.")
		return
	}

	// Check if it's the player's turn
	if playerBattle.CurrentTurn != message.Author.ID {
		session.ChannelMessageSend(message.ChannelID, "It's not your turn!")
		return
	}

	// Find the target (the opponent)
	var targetID string
	for id := range playerBattle.Participants {
		if id != message.Author.ID {
			targetID = id
			break
		}
	}

	// Set the action
	err := playerBattle.SetAction(message.Author.ID, actionName, targetID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Process the turn
	result, err := playerBattle.ProcessTurn()
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error processing turn: %v", err))
		return
	}

	// Update the battle embed
	updateBattleEmbed(session, playerBattle)

	// Check if battle is complete
	if playerBattle.State == BattleComplete {
		handleBattleCompletion(session, playerBattle, battleID)
		return
	}

	// If next turn is a bot/NPC, process it automatically
	if playerBattle.Participants[playerBattle.CurrentTurn].IsBot {
		processBotTurn(session, playerBattle)
	}
}

// processBotTurn automatically processes a turn for an NPC
func processBotTurn(session *discordgo.Session, battle *Battle) {
	// Small delay to make it feel more natural
	time.Sleep(1 * time.Second)

	// Process the turn (NPC action should have been set automatically)
	result, err := battle.ProcessTurn()
	if err != nil {
		fmt.Printf("Error processing NPC turn: %v\n", err)
		return
	}

	// Update the battle embed
	updateBattleEmbed(session, battle)

	// Check if battle is complete
	if battle.State == BattleComplete {
		ActiveBattlesMutex.Lock()
		for id, b := range ActiveBattles {
			if b == battle {
				handleBattleCompletion(session, battle, id)
				break
			}
		}
		ActiveBattlesMutex.Unlock()
	}
}

// updateBattleEmbed updates the battle status embed message
func updateBattleEmbed(session *discordgo.Session, battle *Battle) {
	// Create updated embed
	battleEmbed := createBattleEmbed(battle)

	// Update the message
	_, err := session.ChannelMessageEditEmbed(battle.ChannelID, battle.InteractionMessage, battleEmbed)
	if err != nil {
		fmt.Printf("Error updating battle embed: %v\n", err)
	}
}

// createBattleEmbed creates a rich embed for battle display
func createBattleEmbed(battle *Battle) *discordgo.MessageEmbed {
	// Get both participants for display
	var p1, p2 *CombatParticipant
	for _, p := range battle.Participants {
		if p1 == nil {
			p1 = p
		} else {
			p2 = p
		}
	}

	// Create health bar representations (20 characters wide)
	p1HealthPercent := float64(p1.CurrentHP) / float64(p1.MaxHP)
	p2HealthPercent := float64(p2.CurrentHP) / float64(p2.MaxHP)

	p1HealthBar := createProgressBar(p1HealthPercent, 20)
	p2HealthBar := createProgressBar(p2HealthPercent, 20)

	// Create mana bar representations (10 characters wide)
	p1ManaPercent := float64(p1.CurrentMP) / float64(p1.MaxMP)
	p2ManaPercent := float64(p2.CurrentMP) / float64(p2.MaxMP)

	p1ManaBar := createProgressBar(p1ManaPercent, 10)
	p2ManaBar := createProgressBar(p2ManaPercent, 10)

	// Create embed
	battleEmbed := &discordgo.MessageEmbed{
		Title:       "⚔️ Battle ⚔️",
		Description: fmt.Sprintf("Round %d\n\nCurrent turn: **%s**", battle.Round, battle.Participants[battle.CurrentTurn].UserName),
		Color:       0xFF0000,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: fmt.Sprintf("%s (%s)", p1.UserName, p1.Element),
				Value: fmt.Sprintf("HP: %d/%d %s\nMP: %d/%d %s\nStatus: %s",
					p1.CurrentHP, p1.MaxHP, p1HealthBar,
					p1.CurrentMP, p1.MaxMP, p1ManaBar,
					formatStatusEffects(p1)),
				Inline: true,
			},
			{
				Name: fmt.Sprintf("%s (%s)", p2.UserName, p2.Element),
				Value: fmt.Sprintf("HP: %d/%d %s\nMP: %d/%d %s\nStatus: %s",
					p2.CurrentHP, p2.MaxHP, p2HealthBar,
					p2.CurrentMP, p2.MaxMP, p2ManaBar,
					formatStatusEffects(p2)),
				Inline: true,
			},
			{
				Name:  "Battle Log",
				Value: formatBattleLog(battle),
			},
			{
				Name:  "Commands",
				Value: "• `!cb battle attack` - Physical attack\n• `!cb battle magic` - Magical attack\n• `!cb battle defend` - Increase defense\n• `!cb battle item` - Use healing item\n• `!cb battle forfeit` - Give up",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Battle ID: %s", battle.ID),
		},
	}

	return battleEmbed
}

// createProgressBar generates a text-based progress bar
func createProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	filledWidth := int(percent * float64(width))
	emptyWidth := width - filledWidth

	bar := "["
	for i := 0; i < filledWidth; i++ {
		bar += "█"
	}
	for i := 0; i < emptyWidth; i++ {
		bar += "░"
	}
	bar += "]"

	return bar
}

// formatStatusEffects returns a string representation of status effects
func formatStatusEffects(participant *CombatParticipant) string {
	if len(participant.StatusEffects) == 0 {
		return "None"
	}

	statuses := []string{}
	for effect, turns := range participant.StatusEffects {
		statuses = append(statuses, fmt.Sprintf("%s (%d)", effect, turns))
	}

	return strings.Join(statuses, ", ")
}

// formatBattleLog returns the most recent battle log entries
func formatBattleLog(battle *Battle) string {
	// Show the last 4 log entries or all if fewer
	startIdx := len(battle.Log) - 4
	if startIdx < 0 {
		startIdx = 0
	}

	log := ""
	for i := startIdx; i < len(battle.Log); i++ {
		log += "• " + battle.Log[i] + "\n"
	}

	if log == "" {
		return "The battle has just begun!"
	}

	return log
}

// handleBattleCompletion processes rewards and cleanup when a battle ends
func handleBattleCompletion(session *discordgo.Session, battle *Battle, battleID string) {
	// Get battle results
	result, err := battle.GetResult()
	if err != nil {
		fmt.Printf("Error getting battle result: %v\n", err)
		return
	}

	// Award experience and currency to winner
	db := database.GetDB()

	// Only process rewards for human players (not NPCs)
	if !strings.HasPrefix(result.Winner, "npc_") {
		// Add currency to winner
		_, err = database.AddCurrency(db, result.Winner, result.CurrencyGain)
		if err != nil {
			fmt.Printf("Error adding currency: %v\n", err)
		}

		// Experience points could be added here if implementing a leveling system
	}

	// Create result embed
	resultEmbed := &discordgo.MessageEmbed{
		Title:       "🏆 Battle Complete!",
		Description: fmt.Sprintf("**%s** has won the battle!", battle.Participants[result.Winner].UserName),
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Battle Statistics",
				Value: fmt.Sprintf("Rounds: %d\nRemaining HP: %d\n", result.Rounds, result.WinnerHP),
			},
		},
	}

	// Add rewards section if winner is a player
	if !strings.HasPrefix(result.Winner, "npc_") {
		resultEmbed.Fields = append(resultEmbed.Fields, &discordgo.MessageEmbedField{
			Name:  "Rewards",
			Value: fmt.Sprintf("Experience: %d\nCurrency: %d coins", result.Experience, result.CurrencyGain),
		})
	}

	// Send the result message
	session.ChannelMessageSendEmbed(battle.ChannelID, resultEmbed)

	// Remove battle from active battles
	ActiveBattlesMutex.Lock()
	delete(ActiveBattles, battleID)
	ActiveBattlesMutex.Unlock()
}

// showBattleStatus shows the current battle status
func showBattleStatus(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Find the battle this player is in
	var playerBattle *Battle

	ActiveBattlesMutex.Lock()
	for _, battle := range ActiveBattles {
		if _, exists := battle.Participants[message.Author.ID]; exists {
			playerBattle = battle
			break
		}
	}
	ActiveBattlesMutex.Unlock()

	if playerBattle == nil {
		session.ChannelMessageSend(message.ChannelID, "You're not in a battle! Use `!cb battle start` to begin one.")
		return
	}

	// Update the battle embed
	updateBattleEmbed(session, playerBattle)
}

// forfeitBattle allows a player to give up
func forfeitBattle(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Find the battle this player is in
	var playerBattle *Battle
	var battleID string

	ActiveBattlesMutex.Lock()
	for id, battle := range ActiveBattles {
		if _, exists := battle.Participants[message.Author.ID]; exists {
			playerBattle = battle
			battleID = id
			break
		}
	}

	if playerBattle == nil {
		ActiveBattlesMutex.Unlock()
		session.ChannelMessageSend(message.ChannelID, "You're not in a battle!")
		return
	}

	// Get opponent
	var opponentID string
	for id := range playerBattle.Participants {
		if id != message.Author.ID {
			opponentID = id
			break
		}
	}

	// Set opponent as winner
	playerBattle.Participants[message.Author.ID].CurrentHP = 0
	playerBattle.State = BattleComplete

	// Update the battle log
	playerBattle.Log = append(playerBattle.Log, fmt.Sprintf("%s has forfeited the battle!", playerBattle.Participants[message.Author.ID].UserName))

	// Update the battle embed
	updateBattleEmbed(session, playerBattle)

	// Process battle completion
	handleBattleCompletion(session, playerBattle, battleID)

	// Remove from active battles
	delete(ActiveBattles, battleID)
	ActiveBattlesMutex.Unlock()

	// Send forfeit message
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("**%s** has forfeited the battle!", message.Author.Username))
}
