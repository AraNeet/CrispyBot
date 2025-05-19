package main

import (
	"CrispyBot/bugou"
	"CrispyBot/database"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize the database connection - this creates the singleton instance
	db := database.GetDB()
	defer db.Close()

	fmt.Println("Starting CrispyBot...")

	// Start the shop refresh scheduler
	database.StartShopRefreshScheduler(db)

	// Start the Discord bot
	go bugou.StartBot()

	// Uncomment to start the API server
	// go server.StartServer()

	// Graceful shutdown handling
	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-exit

	fmt.Println("Shutting Down")
}
