package bugouhandlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Print(session)
	if message.Author.ID == session.State.User.ID {
		return
	}
	if message.Content == "Ping" {
		session.ChannelMessageSend(message.ChannelID, "Pong!")
		fmt.Println("Sending Message")
	}
	if message.Content == "Pong" {
		session.ChannelMessageSend(message.ChannelID, "Ping!")
		fmt.Println("Sending Message")
	}
}
