package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/rumblefrog/source-chat-relay/src/server/database"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

func main() {
	// bot.InitBot()

	log.Println("Server is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Received exit signal. Terminating.")

	// bot.RelayBot.Close()

	protocol.NetListener.Close()

	database.DBConnection.Close()
}
