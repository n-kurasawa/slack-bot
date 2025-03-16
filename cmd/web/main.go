package main

import (
	"log"
	"net/http"

	"github.com/n-kurasawa/slack-bot/internal/config"
	"github.com/n-kurasawa/slack-bot/internal/db"
	"github.com/n-kurasawa/slack-bot/internal/web"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// データベース接続の初期化
	store, err := db.NewStore(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	// テーブルの作成（必要に応じて）
	if err := store.CreateTable(); err != nil {
		log.Fatal(err)
	}

	// Webハンドラーの初期化
	handler, err := web.NewHandler(store)
	if err != nil {
		log.Fatal(err)
	}

	// ルーティングの設定
	mux := handler.RegisterRoutes()

	// HTTPサーバーの起動
	webPort := cfg.WebPort

	log.Printf("Web UIサーバーを起動しています。ポート: %s", webPort)
	if err := http.ListenAndServe(":"+webPort, mux); err != nil {
		log.Fatal(err)
	}
}
