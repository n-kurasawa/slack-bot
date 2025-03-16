package main

import (
	"log"
	"net/http"

	"github.com/n-kurasawa/slack-bot/internal/config"
	"github.com/n-kurasawa/slack-bot/internal/db"
	"github.com/n-kurasawa/slack-bot/internal/slack"
	slackapi "github.com/slack-go/slack"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := slackapi.New(cfg.SlackBotToken)

	store, err := db.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	if err := store.CreateTable(); err != nil {
		log.Fatal(err)
	}

	handler := slack.NewHandler(client, store.DB, store, cfg.SlackSigningSecret)

	http.Handle("/slack/events", handler)
	log.Printf("Starting server on %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
