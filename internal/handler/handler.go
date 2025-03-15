package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/slack-go/slack"
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

type ImageStore interface {
	GetImage(db *sql.DB) (*Image, error)
}

type Image struct {
	ID   int64
	Data []byte
}

type Handler struct {
	client     *slack.Client
	db         *sql.DB
	imageStore ImageStore
}

func New(client *slack.Client, database *sql.DB, store ImageStore) *Handler {
	return &Handler{
		client:     client,
		db:         database,
		imageStore: store,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var event SlackEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.handleError(w, fmt.Errorf("JSONのデコードに失敗: %w", err), http.StatusBadRequest)
		return
	}

	if err := h.handleEvent(w, &event); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleEvent(w http.ResponseWriter, event *SlackEvent) error {
	// URL Verification
	if event.Type == "url_verification" {
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(event.Challenge))
		return err
	}

	// メッセージイベントの処理
	if event.Type == "event_callback" && event.Event.Type == "message" {
		return h.handleMessage(event)
	}

	return nil
}

func (h *Handler) handleMessage(event *SlackEvent) error {
	switch event.Event.Text {
	case "hello":
		_, _, err := h.client.PostMessage(
			event.Event.Channel,
			slack.MsgOptionText("world", false),
		)
		if err != nil {
			return fmt.Errorf("メッセージの送信に失敗: %w", err)
		}

	case "image":
		img, err := h.imageStore.GetImage(h.db)
		if err != nil {
			return fmt.Errorf("画像の取得に失敗: %w", err)
		}

		_, err = h.client.UploadFileV2(slack.UploadFileV2Parameters{
			Channel: event.Event.Channel,
			Content: string(img.Data),
		})
		if err != nil {
			return fmt.Errorf("画像の送信に失敗: %w", err)
		}
	}

	return nil
}

func (h *Handler) handleError(w http.ResponseWriter, err error, status int) {
	log.Printf("エラー: %v\n", err)
	w.WriteHeader(status)
}
