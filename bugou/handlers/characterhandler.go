package bugouhandlers

import (
	"CrispyBot/database"
	"CrispyBot/roller"
	"fmt"
	"strings"
	"time"

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
		text := fmt.Errorf("error: %w", err)
		session.ChannelMessageSend(message.ChannelID, text.Error())
		return
	}

	// Create an embed message with the character details
	charEmbed := CreateCharacterEmbed(character, message.Author)
	session.ChannelMessageSendEmbed(message.ChannelID, charEmbed)
}

// HandleDeleteCharacterRequest processes a character deletion request
func HandleDeleteCharacterRequest(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Create a confirmation message with buttons
	confirmEmbed := &discordgo.MessageEmbed{
		Title:       "⚠️ Delete Character",
		Description: "Are you sure you want to delete your character? This action cannot be undone!",
		Color:       0xFF0000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "This confirmation will expire in 60 seconds",
		},
	}

	// Add a unique identifier to track this specific deletion request
	deleteRequestID := fmt.Sprintf("delete_%s_%d", message.Author.ID, time.Now().Unix())

	// Create confirm and cancel buttons
	confirmButton := discordgo.Button{
		Label:    "Yes, Delete Character",
		Style:    discordgo.DangerButton,
		CustomID: fmt.Sprintf("%s_confirm", deleteRequestID),
	}

	cancelButton := discordgo.Button{
		Label:    "Cancel",
		Style:    discordgo.SecondaryButton,
		CustomID: fmt.Sprintf("%s_cancel", deleteRequestID),
	}

	// Create the action row with buttons
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{confirmButton, cancelButton},
	}

	// Send the confirmation message with buttons
	_, err := session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Embed:      confirmEmbed,
		Components: []discordgo.MessageComponent{actionRow},
	})

	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Set up a temporary handler for the button interactions
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Ignore interactions that aren't for this specific deletion request
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		customID := i.MessageComponentData().CustomID
		if !strings.HasPrefix(customID, deleteRequestID) {
			return
		}

		// Check which button was pressed
		if strings.HasSuffix(customID, "_confirm") {
			// Process the deletion
			db := database.GetDB()
			err := database.DeleteCharacter(db, message.Author.ID)

			var responseContent string
			if err != nil {
				responseContent = fmt.Sprintf("Failed to delete character: %v", err)
			} else {
				responseContent = "Your character has been deleted. You can create a new one with `!cb roll`."
			}

			// Respond to the interaction
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content:    responseContent,
					Components: []discordgo.MessageComponent{}, // Remove the buttons
				},
			})
		} else if strings.HasSuffix(customID, "_cancel") {
			// Canceled the deletion
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content:    "Character deletion canceled.",
					Components: []discordgo.MessageComponent{}, // Remove the buttons
				},
			})
		}
	})

	// Set up a timeout to remove the buttons after 60 seconds
	time.AfterFunc(60*time.Second, func() {
		// Edit the message to remove the buttons after timeout
		session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: message.ChannelID,
			ID:      message.ID,
			Embed:   confirmEmbed,
		})
	})
}
