package bugouhandlers

import (
	combathandlers "CrispyBot/bugou/combathandlers"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	prefixCommand       = "!cb"
	helpCommand         = "help"
	rollCommand         = "roll"
	statCommand         = "stats"
	shopCommand         = "shop"
	buyCommand          = "buy"
	walletCommand       = "wallet"
	dailyCommand        = "daily"
	inventoryCommand    = "inventory"
	equipCommand        = "equip"
	unequipCommand      = "unequip"
	rerollCommand       = "reroll"
	rerollStatCommand   = "rerollstat"
	rerollStatusCommand = "rerolls"
	deleteCommand       = "delete"
	battleCommand       = "battle" // Added battle command
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
	case shopCommand:
		HandleShopCommand(session, message)
	case buyCommand:
		HandleBuyCommand(session, message, commandParts)
	case walletCommand:
		HandleWalletCommand(session, message)
	case dailyCommand:
		HandleDailyCommand(session, message)
	case inventoryCommand:
		HandleInventoryCommand(session, message)
	case equipCommand:
		HandleEquipCommand(session, message, commandParts)
	case unequipCommand:
		HandleUnequipCommand(session, message)
	case rerollCommand:
		HandleFullRerollCommand(session, message)
	case rerollStatCommand:
		HandleStatRerollCommand(session, message, commandParts)
	case rerollStatusCommand:
		HandleRerollStatusCommand(session, message)
	case deleteCommand:
		HandleDeleteCharacterRequest(session, message)
	case battleCommand:
		combathandlers.HandleBattleCommand(session, message, commandParts)
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
			{
				Name:  "!cb shop",
				Value: "Browse the item shop",
			},
			{
				Name:  "!cb buy [item number]",
				Value: "Buy an item from the shop",
			},
			{
				Name:  "!cb wallet",
				Value: "Check your currency balance",
			},
			{
				Name:  "!cb daily",
				Value: "Collect your daily currency reward",
			},
			{
				Name:  "!cb inventory",
				Value: "View your inventory of purchased items",
			},
			{
				Name:  "!cb unequip",
				Value: "Unequip your currently equipped item",
			},
			{
				Name:  "!cb equip [item number]",
				Value: "Equip an item from your inventory",
			},
			{
				Name:  "!cb reroll",
				Value: "Reroll your entire character (2 per day)",
			},
			{
				Name:  "!cb rerollstat [stat name]",
				Value: "Reroll a specific stat (1 per day)",
			},
			{
				Name:  "!cb rerolls",
				Value: "Check your remaining rerolls for the day",
			},
			{
				Name:  "!cb delete",
				Value: "Delete your current character (requires confirmation)",
			},
			{
				Name:  "!cb battle start [opponent name/mention] [difficulty]",
				Value: "Start a battle with an NPC or another player",
			},
			{
				Name:  "!cb battle [action]",
				Value: "Battle actions: attack, magic, defend, item, status, forfeit",
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "CrispyBot v1.0",
		},
	}

	session.ChannelMessageSendEmbed(channelID, helpEmbed)
}
