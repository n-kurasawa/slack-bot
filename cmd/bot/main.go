package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"

	"github.com/n-kurasawa/slack-bot/internal/db"
)

type SlackEvent struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
	Event     struct {
		Type    string `json:"type"`
		Text    string `json:"text"`
		Channel string `json:"channel"`
		User    string `json:"user"`
	} `json:"event"`
}

func main() {
	// .env ファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf(".env ファイルの読み込みに失敗しました: %v\n", err)
	}

	token := os.Getenv("SLACK_BOT_TOKEN")
	if token == "" {
		log.Fatal("SLACK_BOT_TOKEN を設定してください")
	}

	database, err := db.OpenDB("images.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	client := slack.New(token)

	http.HandleFunc("/slack/events", func(w http.ResponseWriter, r *http.Request) {
		var event SlackEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			log.Printf("JSONのデコードに失敗しました: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// URL Verification
		if event.Type == "url_verification" {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(event.Challenge))
			return
		}

		// メッセージイベントの処理
		if event.Type == "event_callback" && event.Event.Type == "message" {
			switch event.Event.Text {
			case "hello":
				_, _, err := client.PostMessage(event.Event.Channel, slack.MsgOptionText("world", false))
				if err != nil {
					log.Printf("メッセージの送信に失敗しました: %v\n", err)
				}
			case "image":
				img, err := db.GetImage(database)
				if err != nil {
					log.Printf("画像の取得に失敗しました: %v\n", err)
					return
				}

				_, err = client.UploadFile(slack.FileUploadParameters{
					Channels: []string{event.Event.Channel},
					Content:  string(img.Data),
				})
				if err != nil {
					log.Printf("画像の送信に失敗しました: %v\n", err)
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	})

	port := "3000"
	fmt.Printf("Slack Bot Starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
