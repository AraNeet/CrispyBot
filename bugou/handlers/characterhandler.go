package bugouhandlers

import (
	"CrispyBot/database"
	"CrispyBot/roller"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// handleRollCommand generates a new character for the user
func HandleRollCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
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
	charEmbed := CreateCharacterEmbed(savedChar, message.Author)
	session.ChannelMessageSendEmbed(message.ChannelID, charEmbed)
}

// handleStatsCommand shows the user's character stats
func HandleStatsCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
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
	charEmbed := CreateCharacterEmbed(character, message.Author)
	session.ChannelMessageSendEmbed(message.ChannelID, charEmbed)
}
