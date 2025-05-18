package bugou

import (
	bugouhandlers "CrispyBot/bugou/handlers"
	"CrispyBot/variables"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func StartBot() {
	session, err := discordgo.New("Bot " + variables.Bottoken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	session.AddHandler(bugouhandlers.MessageCreate)
	session.Identify.Intents = discordgo.IntentGuildMessages

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	select {}
}
