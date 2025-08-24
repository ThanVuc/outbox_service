package main

import (
	"log"
	"outbox_service/cmd"
)

func main() {
	log.Println("Outbox service start listening")
	cmd.OutboxRun()
}
