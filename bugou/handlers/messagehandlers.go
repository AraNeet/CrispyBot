package bugouhandlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	prefixCommand = "!cb"
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
		SendHelpMessage(session, message.ChannelID)
		return
	}

	command := strings.ToLower(commandParts[1])

	// Handle commands
	switch command {
	case helpCommand:
		SendHelpMessage(session, message.ChannelID)
	case rollCommand:
		HandleRollCommand(session, message)
	case statCommand:
		HandleStatsCommand(session, message)
	default:
		session.ChannelMessageSend(message.ChannelID, "Unknown command. Try `!cb help` for a list of commands.")
	}
}

// sendHelpMessage sends the help message with available commands
func SendHelpMessage(session *discordgo.Session, channelID string) {
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
