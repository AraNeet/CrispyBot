package bugouhandlers

import (
	"CrispyBot/database"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// HandleInventoryCommand displays the user's inventory
func HandleInventoryCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.DBInit()

	// Get user info
	user, err := database.GetUserByID(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Create an inventory embed
	inventoryEmbed := &discordgo.MessageEmbed{
		Title:       "🎒 Your Inventory",
		Description: fmt.Sprintf("%s's collection of items", message.Author.Username),
		Color:       0x964B00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use !cb shop to find more items",
		},
	}

	// Add inventory items to the embed
	if len(user.Inventory) == 0 {
		inventoryEmbed.Fields = append(inventoryEmbed.Fields, &discordgo.MessageEmbedField{
			Name:  "Empty Inventory",
			Value: "You don't have any items yet. Visit the shop with `!cb shop` to buy some!",
		})
	} else {
		// Group items by type for cleaner display
		weaponsList := ""
		itemCount := 0

		for _, itemName := range user.Inventory {
			itemCount++
			weaponsList += fmt.Sprintf("• %s\n", itemName)
		}

		inventoryEmbed.Fields = append(inventoryEmbed.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("Weapons (%d)", itemCount),
			Value: weaponsList,
		})
	}

	session.ChannelMessageSendEmbed(message.ChannelID, inventoryEmbed)
}
