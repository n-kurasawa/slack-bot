package main

import (
	"log"
	"net/http"

	"github.com/n-kurasawa/slack-bot/internal/bot"
	"github.com/n-kurasawa/slack-bot/internal/config"
	"github.com/n-kurasawa/slack-bot/internal/db"
	"github.com/slack-go/slack"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := slack.New(cfg.SlackBotToken)
	store, err := db.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	handler := bot.NewHandler(client, store.DB, store, cfg.SlackSigningSecret)
	http.Handle("/slack/events", handler)

	log.Printf("Starting server on port %s", cfg.BotPort)
	if err := http.ListenAndServe(":"+cfg.BotPort, nil); err != nil {
		log.Fatal(err)
	}
}
