package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var (
	token string = os.Getenv("BOTTOKEN")
)

func main() {
	BotSession, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Fail to create session.")
	}
	defer BotSession.Close()

	BotSession.AddHandler(messageCreate)
	BotSession.Identify.Intents = 1 << 9

	err = BotSession.Open()
	if err != nil {
		fmt.Println("Error connecting. ", err)
		return
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-exit
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if message.Content == "Ping" {
		session.ChannelMessageSend(message.ChannelID, "Pong!")
	}
	if message.Content == "Pong" {
		session.ChannelMessageSend(message.ChannelID, "Ping!")
	}
}
