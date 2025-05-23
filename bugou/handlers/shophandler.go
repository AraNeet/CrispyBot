package bugouhandlers

import (
	"CrispyBot/database"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// HandleShopCommand displays the current shop inventory
func HandleShopCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.DBInit()

	// Get the current shop
	shop, err := database.GetShop(db)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error accessing the shop: %v", err))
		return
	}

	// Check if user has a wallet, if not initialize it with 500 coins
	err = database.InitializeUserWallet(db, message.Author.ID, 500)
	if err != nil {
		fmt.Printf("Error initializing wallet: %v\n", err)
	}

	// Create an embed message with the shop details
	shopEmbed := &discordgo.MessageEmbed{
		Title:       "🛒 Item Shop",
		Description: fmt.Sprintf("The shop will refresh in %s", formatDuration(time.Until(shop.Timer))),
		Color:       0x00AAFF,
		Fields:      []*discordgo.MessageEmbedField{},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use !cb buy [number] to purchase an item",
		},
	}

	// Add items to the embed
	if len(shop.Inventory.Items) == 0 {
		shopEmbed.Fields = append(shopEmbed.Fields, &discordgo.MessageEmbedField{
			Name:  "No Items Available",
			Value: "The shop is currently empty. Check back after it refreshes.",
		})
	} else {
		for idx, item := range shop.Inventory.Items {
			// Format item stats
			statsText := formatItemStats(item.Stats)

			itemField := &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%d. %s (%s) - %d coins", idx, item.Name, item.Rarity, item.Price),
				Value: statsText,
			}
			shopEmbed.Fields = append(shopEmbed.Fields, itemField)
		}
	}

	session.ChannelMessageSendEmbed(message.ChannelID, shopEmbed)
}

// HandleBuyCommand processes a purchase from the shop
func HandleBuyCommand(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	// Check if the user provided an item number
	if len(args) < 3 {
		session.ChannelMessageSend(message.ChannelID, "Please specify an item number to buy. Usage: `!cb buy [number]`")
		return
	}

	// Parse the item number
	itemIdx, err := strconv.Atoi(args[2])
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Invalid item number. Please provide a valid number.")
		return
	}

	// Get the database singleton
	db := database.DBInit()

	// Process the purchase
	item, err := database.BuyItem(db, message.Author.ID, itemIdx)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Purchase failed: %v", err))
		return
	}

	// Format item stats
	statsText := formatItemStats(item.Stats)

	// Create a purchase confirmation embed
	purchaseEmbed := &discordgo.MessageEmbed{
		Title:       "Purchase Successful",
		Description: fmt.Sprintf("You bought **%s** for **%d** coins!", item.Name, item.Price),
		Color:       0x00FF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Item Details",
				Value: fmt.Sprintf("**Rarity:** %s\n%s", item.Rarity, statsText),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Type !cb wallet to check your remaining balance",
		},
	}

	session.ChannelMessageSendEmbed(message.ChannelID, purchaseEmbed)
}

// HandleWalletCommand shows a user's currency balance
func HandleWalletCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.DBInit()

	// Initialize wallet if needed
	err := database.InitializeUserWallet(db, message.Author.ID, 500)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Get user info
	user, err := database.GetUserByID(db, message.Author.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Count inventory items
	inventoryCount := 0
	if user.Inventory != nil {
		inventoryCount = len(user.Inventory)
	}

	// Create a wallet embed
	walletEmbed := &discordgo.MessageEmbed{
		Title:       "💰 Your Wallet",
		Description: fmt.Sprintf("**Balance:** %d coins", user.Wallet),
		Color:       0xFFD700,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Inventory",
				Value: fmt.Sprintf("You have %d items in your inventory", inventoryCount),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use !cb daily to earn daily rewards",
		},
	}

	session.ChannelMessageSendEmbed(message.ChannelID, walletEmbed)
}

// HandleDailyCommand gives the user their daily currency reward
func HandleDailyCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the database singleton
	db := database.DBInit()

	// TODO: Implement daily reward cooldown
	// For now, just give 100 coins every time
	newBalance, err := database.AddCurrency(db, message.Author.ID, 100)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Create a reward embed
	rewardEmbed := &discordgo.MessageEmbed{
		Title:       "Daily Reward",
		Description: "You received **100** coins!",
		Color:       0xFFD700,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "New Balance",
				Value: fmt.Sprintf("%d coins", newBalance),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Come back tomorrow for more rewards",
		},
	}

	session.ChannelMessageSendEmbed(message.ChannelID, rewardEmbed)
}

// Helper function to format item stats
func formatItemStats(stats map[string]int) string {
	if len(stats) == 0 {
		return "No stat bonuses"
	}

	var statsText string
	for stat, value := range stats {
		if value > 0 {
			statsText += fmt.Sprintf("• +%d to %s\n", value, stat)
		} else {
			statsText += fmt.Sprintf("• %d to %s\n", value, stat)
		}
	}

	return statsText
}

// Helper function to format duration until shop refresh
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d hours and %d minutes", hours, minutes)
	}
	return fmt.Sprintf("%d minutes", minutes)
}
