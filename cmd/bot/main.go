package main

import (
	"log"
	"net/http"

	"github.com/n-kurasawa/slack-bot/internal/bot"
	"github.com/n-kurasawa/slack-bot/internal/db"
	slackapi "github.com/slack-go/slack"
)

type Config struct {
	SlackBotToken string
	DBPath        string
	Port          string
}

func main() {
	cfg := &Config{
		SlackBotToken: "your-token",
		DBPath:        "images.db",
		Port:          "8080",
	}

	client := slackapi.New(cfg.SlackBotToken)
	store, err := db.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	handler := bot.NewHandler(client, store.DB, store, "your-signing-secret")
	http.Handle("/slack/events", handler)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
