package bugouhandlers

import (
	"CrispyBot/database"
	"CrispyBot/database/models"
	"CrispyBot/roller"
	"CrispyBot/variables"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// HandleFullRerollCommand completely rerolls a user's character
func HandleFullRerollCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.DBInit()

	// Use one full reroll
	remainingRerolls, err := database.UseFullReroll(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Reroll failed: %v", err))
		return
	}

	// Delete existing character (if any)
	_, err = database.GetCharacterByOwner(db, message.Author.ID)
	if err == nil {
		database.DeleteCharacter(db, message.Author.ID)
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
	charEmbed := CreateCharacterEmbed(savedChar, message.Author)

	// Add reroll info to the footer
	charEmbed.Footer.Text = fmt.Sprintf("%s | Remaining full rerolls today: %d", charEmbed.Footer.Text, remainingRerolls)

	session.ChannelMessageSendEmbed(message.ChannelID, charEmbed)
}

// HandleStatRerollCommand rerolls a single stat
func HandleStatRerollCommand(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	// Check if stat type was specified
	if len(args) < 3 {
		session.ChannelMessageSend(message.ChannelID, "Please specify which stat to reroll. Usage: `!cb rerollstat [vitality|strength|speed|durability|intelligence|mana|mastery]`")
		return
	}

	// Parse the stat type
	statArg := strings.ToLower(args[2])
	var statType variables.StatType

	switch statArg {
	case "vitality":
		statType = variables.Vitality
	case "strength":
		statType = variables.Strength
	case "speed":
		statType = variables.Speed
	case "durability":
		statType = variables.Durability
	case "intelligence":
		statType = variables.Intelligence
	case "mana":
		statType = variables.Mana
	case "mastery":
		statType = variables.Mastery
	default:
		session.ChannelMessageSend(message.ChannelID, "Invalid stat type. Valid stats are: vitality, strength, speed, durability, intelligence, mana, mastery")
		return
	}

	// Get the database singleton
	db := database.DBInit()

	// Use a stat reroll
	remainingRerolls, err := database.UseStatReroll(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Stat reroll failed: %v", err))
		return
	}

	// Get the character before reroll for comparison
	oldCharacter, err := database.GetCharacterByOwner(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Get current stat value for comparison
	var oldStat models.Stat
	var oldStatName string
	switch statType {
	case variables.Vitality:
		oldStat = oldCharacter.Stats.Vitality
		oldStatName = "Vitality"
	case variables.Strength:
		oldStat = oldCharacter.Stats.Strength
		oldStatName = "Strength"
	case variables.Speed:
		oldStat = oldCharacter.Stats.Speed
		oldStatName = "Speed"
	case variables.Durability:
		oldStat = oldCharacter.Stats.Durability
		oldStatName = "Durability"
	case variables.Intelligence:
		oldStat = oldCharacter.Stats.Intelligence
		oldStatName = "Intelligence"
	case variables.Mana:
		oldStat = oldCharacter.Stats.Mana
		oldStatName = "Mana"
	case variables.Mastery:
		oldStat = oldCharacter.Stats.Mastery
		oldStatName = "Mastery"
	}

	// Reroll the stat
	newStat, err := database.RerollSingleStat(db, message.Author.ID, statType)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Failed to reroll stat: %v", err))
		return
	}

	// Get the updated character
	updatedChar, err := database.GetCharacterByOwner(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Create a stat reroll embed
	rerollEmbed := &discordgo.MessageEmbed{
		Title:       "Stat Reroll Result",
		Description: fmt.Sprintf("You rerolled your %s stat!", oldStatName),
		Color:       0x00AAFF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Before",
				Value:  fmt.Sprintf("%s: %d (%s) [%s]", oldStatName, oldStat.Value, oldStat.Stat_Name, oldStat.Rarity),
				Inline: true,
			},
			{
				Name:   "After",
				Value:  fmt.Sprintf("%s: %d (%s) [%s]", oldStatName, newStat.Value, newStat.Stat_Name, newStat.Rarity),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Remaining stat rerolls today: %d", remainingRerolls),
		},
	}

	session.ChannelMessageSendEmbed(message.ChannelID, rerollEmbed)

	// Also send the updated character sheet
	charEmbed := CreateCharacterEmbed(updatedChar, message.Author)
	session.ChannelMessageSendEmbed(message.ChannelID, charEmbed)
}

// HandleRerollStatusCommand shows remaining rerolls for the day
func HandleRerollStatusCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.DBInit()

	// Get the user
	user, err := database.GetUserByID(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Get the shop (for timer)
	shop, err := database.GetShop(db)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Create reroll status embed
	rerollEmbed := &discordgo.MessageEmbed{
		Title:       "🎲 Reroll Status",
		Description: "Your daily reroll allowance",
		Color:       0x9B59B6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Full Character Rerolls",
				Value:  fmt.Sprintf("%d/2 remaining", user.FullRerolls),
				Inline: true,
			},
			{
				Name:   "Single Stat Rerolls",
				Value:  fmt.Sprintf("%d/1 remaining", user.StatRerolls),
				Inline: true,
			},
			{
				Name:  "Refresh Time",
				Value: fmt.Sprintf("Rerolls refresh in %s", formatDuration(time.Until(shop.Timer))),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use !cb reroll to reroll your entire character or !cb rerollstat [stat] to reroll a specific stat",
		},
	}

	session.ChannelMessageSendEmbed(message.ChannelID, rerollEmbed)
}
