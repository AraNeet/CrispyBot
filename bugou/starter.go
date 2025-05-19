package bugou

import (
	bugouhandlers "CrispyBot/bugou/handlers"
	"CrispyBot/variables"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func StartBot() {
	// Initialize Discord bot
	session, err := discordgo.New("Bot " + variables.Bottoken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	// Register message handler
	session.AddHandler(bugouhandlers.MessageCreate)
	session.Identify.Intents = discordgo.IntentGuildMessages

	// Open websocket connection to Discord
	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}
	defer session.Close()

	fmt.Println("Discord bot is now running. Press CTRL+C to exit.")
	select {}
}
