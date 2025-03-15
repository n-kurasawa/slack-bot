package main

import (
	"log"
	"net/http"

	slackapi "github.com/slack-go/slack"

	"github.com/n-kurasawa/slack-bot/internal/config"
	"github.com/n-kurasawa/slack-bot/internal/image"
	"github.com/n-kurasawa/slack-bot/internal/slack"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := slackapi.New(cfg.SlackBotToken)

	store, err := image.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	if err := store.CreateTable(); err != nil {
		log.Fatal(err)
	}

	handler := slack.NewHandler(client, store.DB, store)

	http.Handle("/slack/events", handler)
	log.Printf("Starting server on %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
