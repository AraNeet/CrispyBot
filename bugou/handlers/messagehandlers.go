package bugouhandlers

import (
	"CrispyBot/database"
	"CrispyBot/database/models"
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
func createCharacterEmbed(character models.Character, author *discordgo.User) *discordgo.MessageEmbed {
	// Get the character characteristics, traits, and stats
	chars := character.Characteristics
	stats := character.Stats
	traits := character.Attributes

	// Create the embed
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s's Character", author.Username),
		Description: fmt.Sprintf("Race: **%s** | Element: **%s** | Alignment: **%s**", chars.Race.Trait_Name, chars.Element.Trait_Name, chars.Alignment.Trait_Name),
		Color:       0xFF5500,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: author.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Characteristics",
				Value:  formatCharacteristics(chars),
				Inline: false,
			},
			{
				Name:   "Stats",
				Value:  formatStats(stats),
				Inline: true,
			},
			{
				Name:   "Traits",
				Value:  formatTraits(traits),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Character ID: %s", character.ID.Hex()),
		},
	}

	return embed
}

// formatCharacteristics formats the character characteristics for display
func formatCharacteristics(chars models.Characteristics) string {
	var charDetails string

	// Format Race with its rarity
	charDetails += fmt.Sprintf("**Race:** %s (%s)\n", chars.Race.Trait_Name, chars.Race.Rarity)

	// Format Element with its rarity
	charDetails += fmt.Sprintf("**Element:** %s (%s)\n", chars.Element.Trait_Name, chars.Element.Rarity)

	// Format Alignment with its rarity
	charDetails += fmt.Sprintf("**Alignment:** %s (%s)\n", chars.Alignment.Trait_Name, chars.Alignment.Rarity)

	// Add any race-specific stat bonuses or penalties
	if len(chars.Race.Stats_Value) > 0 {
		charDetails += "\n**Race Bonuses/Penalties:**\n"
		for stat, value := range chars.Race.Stats_Value {
			if value > 0 {
				charDetails += fmt.Sprintf("• +%d to %s\n", value, stat)
			} else if value < 0 {
				charDetails += fmt.Sprintf("• %d to %s\n", value, stat)
			}
		}
	}

	return charDetails
}

// formatStats formats the character stats for display
func formatStats(stats models.StatsSheets) string {
	// Create a uniform format for all stats with name, value, and rarity
	return fmt.Sprintf(
		"**Vitality:** %d (%s) [%s]\n**Strength:** %d (%s) [%s]\n**Speed:** %d (%s) [%s]\n**Durability:** %d (%s) [%s]\n**Intelligence:** %d (%s) [%s]\n**Mana:** %d (%s) [%s]\n**Mastery:** %d (%s) [%s]",
		stats.Vitality.Value, stats.Vitality.Stat_Name, stats.Vitality.Rarity,
		stats.Strength.Value, stats.Strength.Stat_Name, stats.Strength.Rarity,
		stats.Speed.Value, stats.Speed.Stat_Name, stats.Speed.Rarity,
		stats.Durability.Value, stats.Durability.Stat_Name, stats.Durability.Rarity,
		stats.Intelligence.Value, stats.Intelligence.Stat_Name, stats.Intelligence.Rarity,
		stats.Mana.Value, stats.Mana.Stat_Name, stats.Mana.Rarity,
		stats.Mastery.Value, stats.Mastery.Stat_Name, stats.Mastery.Rarity,
	)
}

// formatTraits formats the character traits for display
func formatTraits(traits models.Traits) string {
	var traitDetails string

	// Format Innate trait (if not "None")
	if traits.Innate.Trait_Name != "None" {
		traitDetails += fmt.Sprintf("**Innate Trait:** %s (%s)\n", traits.Innate.Trait_Name, traits.Innate.Rarity)

		// Show stat bonuses from innate trait if any
		if len(traits.Innate.Stats_Value) > 0 {
			traitDetails += "**Bonuses:**\n"
			for stat, value := range traits.Innate.Stats_Value {
				if value != 0 {
					traitDetails += fmt.Sprintf("• %+d to %s\n", value, stat)
				}
			}
		}
	}

	// Format Inadequacy trait (if not "None")
	if traits.Inadequacy.Trait_Name != "None" {
		traitDetails += fmt.Sprintf("\n**Weakness:** %s\n", traits.Inadequacy.Trait_Name)

		// Show stat penalties from inadequacy trait if any
		if len(traits.Inadequacy.Stats_Value) > 0 {
			traitDetails += "**Penalties:**\n"
			for stat, value := range traits.Inadequacy.Stats_Value {
				if value != 0 {
					traitDetails += fmt.Sprintf("• %+d to %s\n", value, stat)
				}
			}
		}
	}

	// Format X-Factor trait (if not "None")
	if traits.X_Factor.Trait_Name != "None" {
		traitDetails += fmt.Sprintf("\n**X-Factor:** %s\n", traits.X_Factor.Trait_Name)

		// Show stat modifiers from X-Factor trait if any
		if len(traits.X_Factor.Stats_Value) > 0 {
			traitDetails += "**Effects:**\n"
			for stat, value := range traits.X_Factor.Stats_Value {
				if value != 0 {
					traitDetails += fmt.Sprintf("• %+d to %s\n", value, stat)
				}
			}
		}
	}

	// If there are no traits, provide a message
	if traitDetails == "" {
		traitDetails = "No special traits"
	}

	return traitDetails
}
