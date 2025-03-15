package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/slack-go/slack"

	"github.com/n-kurasawa/slack-bot/internal/config"
	"github.com/n-kurasawa/slack-bot/internal/db"
	"github.com/n-kurasawa/slack-bot/internal/handler"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// データベースの初期化
	store, err := db.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if err := store.CreateTable(); err != nil {
		log.Fatal(err)
	}

	// Slackクライアントの初期化
	client := slack.New(cfg.SlackBotToken)

	// ハンドラーの初期化
	handler := handler.New(client, store.DB, store)

	// サーバーの起動
	http.Handle("/slack/events", handler)

	addr := ":" + cfg.Port
	fmt.Printf("Slack Bot Starting on port %s...\n", cfg.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
