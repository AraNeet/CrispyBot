package main

import (
	"CrispyBot/bugou"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go bugou.StartBot()

	// go server.StartServer()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-exit
	fmt.Println("Shutting Down")
}
