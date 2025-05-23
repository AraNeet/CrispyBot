package bugouhandlers

import (
	"CrispyBot/database"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// HandleEquipCommand equips an item from the user's inventory
func HandleEquipCommand(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	// Check if the user provided an item key
	if len(args) < 3 {
		session.ChannelMessageSend(message.ChannelID, "Please specify which item to equip. Usage: `!cb equip [item number]`")
		return
	}

	// Get the item key
	itemKey := "weapon_" + args[2]

	// Get the database singleton
	db := database.DBInit()

	// Get user info
	user, err := database.GetUserByID(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Check if item exists in inventory
	itemName, ok := user.Inventory[itemKey]
	if !ok {
		session.ChannelMessageSend(message.ChannelID, "Item not found in your inventory.")
		return
	}

	// Equip the item
	err = database.EquipItem(db, message.Author.ID, itemKey)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Failed to equip item: %v", err))
		return
	}

	// Get the item stats
	item, err := database.GetItem(db, message.Author.ID, itemKey)
	if err != nil {
		// If we can't find detailed stats, just show success message
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Successfully equipped **%s**!", itemName))
		return
	}

	// Format item stats
	statsText := formatItemStats(item.Stats)

	// Create an equip confirmation embed
	equipEmbed := &discordgo.MessageEmbed{
		Title:       "Item Equipped",
		Description: fmt.Sprintf("You equipped **%s**!", itemName),
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Item Details",
				Value: fmt.Sprintf("**Rarity:** %s\n%s", item.Rarity, statsText),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Type !cb stats to see your updated character stats",
		},
	}

	session.ChannelMessageSendEmbed(message.ChannelID, equipEmbed)
}

// HandleUnequipCommand removes the currently equipped item
func HandleUnequipCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.DBInit()

	// Get character info to check if there's an equipped weapon
	character, err := database.GetCharacterByOwner(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Check if there's anything equipped
	if character.EquippedWeapon.ItemKey == "" {
		session.ChannelMessageSend(message.ChannelID, "You don't have any item equipped.")
		return
	}

	// Store the item name before unequipping
	equippedItemName := character.EquippedWeapon.ItemName

	// Unequip the item
	err = database.UnequipItem(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Failed to unequip item: %v", err))
		return
	}

	// Create an unequip confirmation embed
	unequipEmbed := &discordgo.MessageEmbed{
		Title:       "Item Unequipped",
		Description: fmt.Sprintf("You unequipped **%s**.", equippedItemName),
		Color:       0x00AAFF,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Type !cb stats to see your updated character stats",
		},
	}

	session.ChannelMessageSendEmbed(message.ChannelID, unequipEmbed)
}
