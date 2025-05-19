package bugouhandlers

import (
	"CrispyBot/database"
	"CrispyBot/roller"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	prefixCommand = "!cb"
	initCommand   = "init"
	testCommand   = "test"
	helpCommand   = "help"
	rollCommand   = "roll"
	statCommand   = "stats"
)

// MessageCreate handles incoming Discord messages
func MessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if message.Author.ID == session.State.User.ID {
		return
	}

	// Check if message starts with the command prefix
	if !strings.HasPrefix(message.Content, prefixCommand) {
		return
	}

	// Parse the command
	commandParts := strings.Fields(message.Content)
	if len(commandParts) < 2 {
		// Just the prefix, show help message
		sendHelpMessage(session, message.ChannelID)
		return
	}

	command := strings.ToLower(commandParts[1])

	// Handle commands
	switch command {
	case testCommand:
		handleTestsCommand(session, message)
	case helpCommand:
		sendHelpMessage(session, message.ChannelID)
	case rollCommand:
		handleRollCommand(session, message)
	case statCommand:
		handleStatsCommand(session, message)
	case initCommand:
		handleInitCommand(session, message)
	default:
		session.ChannelMessageSend(message.ChannelID, "Unknown command. Try `!cb help` for a list of commands.")
	}
}

func handleTestsCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.GetDB()

	user, err := database.CreateUser(db, message.Author.ID)
	if err != nil {
		errEmbed := &discordgo.MessageEmbed{
			Title:       "Error message",
			Description: err.Error(),
			Color:       0x00AAFF,
		}
		session.ChannelMessageSendEmbed(message.ChannelID, errEmbed)
		return
	}
	sucEmbed := &discordgo.MessageEmbed{
		Title:       "Create message",
		Description: user.DiscordID,
		Color:       0x00AAFF,
	}
	session.ChannelMessageSendEmbed(message.ChannelID, sucEmbed)
}

func handleInitCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Implementation for init command
}

// sendHelpMessage sends the help message with available commands
func sendHelpMessage(session *discordgo.Session, channelID string) {
	helpEmbed := &discordgo.MessageEmbed{
		Title:       "CrispyBot Help",
		Description: "Here are the commands you can use:",
		Color:       0x00AAFF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "!cb help",
				Value: "Shows this help message",
			},
			{
				Name:  "!cb roll",
				Value: "Rolls a new character for you",
			},
			{
				Name:  "!cb stats",
				Value: "Shows your character's stats",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "CrispyBot v1.0",
		},
	}

	session.ChannelMessageSendEmbed(channelID, helpEmbed)
}

// handleRollCommand generates a new character for the user
func handleRollCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.GetDB()

	// Check if user already has a character
	existingChar, err := database.GetCharacterByOwner(db, message.Author.ID)
	if err == nil {
		// User already has a character
		characterID := existingChar.ID.Hex()
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("You already have a character with ID **%s**! Use `!cb stats` to see your character.", characterID))
		return
	}

	// Generate a new character
	newCharacter := roller.GenerateCharacter(message.Author.ID)

	// Save to the database
	savedChar, err := database.SaveCharacter(db, newCharacter, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error generating character: %v", err))
		return
	}

	// Create an embed message with the character details
	charEmbed := createCharacterEmbed(savedChar, message.Author)
	session.ChannelMessageSendEmbed(message.ChannelID, charEmbed)
}

// handleStatsCommand shows the user's character stats
func handleStatsCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.GetDB()

	// Get the user's character
	character, err := database.GetCharacterByOwner(db, message.Author.ID)
	if err != nil {
		text := fmt.Errorf("Error: %w", err)
		session.ChannelMessageSend(message.ChannelID, text.Error())
		return
	}

	// Create an embed message with the character details
	charEmbed := createCharacterEmbed(character, message.Author)
	session.ChannelMessageSendEmbed(message.ChannelID, charEmbed)
}

// createCharacterEmbed creates a Discord embed message with character details
func createCharacterEmbed(character database.Character, author *discordgo.User) *discordgo.MessageEmbed {
	// Get the character attributes
	attrs := character.Attributes
	stats := character.Stats

	// Create the embed
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s's Character", author.Username),
		Description: fmt.Sprintf("Race: **%s** | Element: **%s** | Alignment: **%s**", attrs.Race.Trait_Name, attrs.Element.Trait_Name, attrs.Alignment.Trait_Name),
		Color:       0xFF5500,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: author.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Stats",
				Value:  formatStats(stats),
				Inline: true,
			},
			{
				Name:   "Traits",
				Value:  formatTraits(attrs),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Character ID: %s", character.ID.Hex()),
		},
	}

	return embed
}

// formatStats formats the character stats for display
func formatStats(stats database.StatsSheets) string {
	return fmt.Sprintf(
		"**Vitality:** %d (%s)\n**Strength:** %d (%s)\n**Speed:** %d (%s)\n**Durability:** %d (%s)\n**Intelligence:** %d (%s)\n**Mana Flow:** %d (%s)\n**Skill Level:** %d (%s)",
		stats.Vitality.Value, stats.Vitality.Stat_Name,
		stats.Strength.Value, stats.Strength.Stat_Name,
		stats.Speed.Value, stats.Speed.Stat_Name,
		stats.Durability.Value, stats.Durability.Stat_Name,
		stats.Intelligence.Value, stats.Intelligence.Stat_Name,
		stats.Mana.Value, stats.Mana.Stat_Name,
		stats.Mastery.Value, stats.Mastery.Stat_Name,
	)
}

// formatTraits formats the character traits for display
func formatTraits(attrs database.Attributes) string {
	var traitDetails string

	traitDetails += fmt.Sprintf("**Race:** %s (%s)\n", attrs.Race.Trait_Name, attrs.Race.Rarity)
	traitDetails += fmt.Sprintf("**Element:** %s\n", attrs.Element.Trait_Name)

	if attrs.Trait.Trait_Name != "None" {
		traitDetails += fmt.Sprintf("**Trait:** %s (%s)\n", attrs.Trait.Trait_Name, attrs.Trait.Rarity)
	}

	if attrs.Weakness.Trait_Name != "None" {
		traitDetails += fmt.Sprintf("**Weakness:** %s\n", attrs.Weakness.Trait_Name)
	}

	if attrs.X_Factor.Trait_Name != "None" {
		traitDetails += fmt.Sprintf("**X-Factor:** %s\n", attrs.X_Factor.Trait_Name)
	}

	return traitDetails
}
