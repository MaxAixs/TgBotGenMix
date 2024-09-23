package main

import (
	"BotMixology/client/telegram"
	event_consumer "BotMixology/consumer/event-consumer"
	telegram2 "BotMixology/events/telegram"
	"BotMixology/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.NewClient(tgBotHost, mustToken())
	store := files.NewStorage()
	processor := telegram2.NewProcessor(tgClient, store)
	consumer := event_consumer.NewConsumer(processor, processor, 100)

	log.Println("Service starting")

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped: ", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access telegram bot!",
	)

	flag.Parse()
	if *token == "" {
		log.Fatal("token is empty")
	}

	return *token
}
