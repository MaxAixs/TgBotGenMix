package main

import (
	"BotMixology/client/telegram"
	event_consumer "BotMixology/consumer/event-consumer"
	telegram2 "BotMixology/events/telegram"
	"BotMixology/storage/sqlite"
	"flag"
	_ "github.com/glebarez/go-sqlite"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "D:/SQLite/my_database.db"
)

func main() {
	tgClient := telegram.NewClient(tgBotHost, mustToken())
	store, err := sqlite.NewDb(sqliteStoragePath)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	processor := telegram2.NewProcessor(tgClient, *store)
	consumer := event_consumer.NewConsumer(processor, processor, 100)

	log.Println("Service starting")

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped: ", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"7380869857:AAFRlQCT5qVxqpubXKgUFYjXqzEbz9k0uHo",
		"token for access telegram bot!",
	)

	flag.Parse()
	if *token == "" {
		log.Fatal("token is empty")
	}

	return *token
}
