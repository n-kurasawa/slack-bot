package slack

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/n-kurasawa/slack-bot/internal/image"
	slackapi "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type SlackClient interface {
	PostMessage(channelID string, options ...slackapi.MsgOption) (string, string, error)
	UploadFileV2(params slackapi.UploadFileV2Parameters) (*slackapi.FileSummary, error)
}

type ImageStore interface {
	GetImage() (*image.Image, error)
}

type Handler struct {
	client     SlackClient
	db         *sql.DB
	imgStore   ImageStore
	signingKey string
}

func NewHandler(client SlackClient, database *sql.DB, store ImageStore, signingKey string) *Handler {
	return &Handler{
		client:     client,
		db:         database,
		imgStore:   store,
		signingKey: signingKey,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleError(w, fmt.Errorf("リクエストボディの読み取りに失敗: %w", err), http.StatusBadRequest)
		return
	}

	if err := h.verifySignature(r.Header, body); err != nil {
		h.handleError(w, err, http.StatusUnauthorized)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		h.handleError(w, fmt.Errorf("イベントのパースに失敗: %w", err), http.StatusBadRequest)
		return
	}

	if err := h.handleEvent(w, &eventsAPIEvent, body); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) verifySignature(header http.Header, body []byte) error {
	sv, err := slackapi.NewSecretsVerifier(header, h.signingKey)
	if err != nil {
		return fmt.Errorf("署名検証の初期化に失敗: %w", err)
	}
	if _, err := sv.Write(body); err != nil {
		return fmt.Errorf("署名の検証に失敗: %w", err)
	}
	if err := sv.Ensure(); err != nil {
		return fmt.Errorf("署名が無効です: %w", err)
	}
	return nil
}

func (h *Handler) handleEvent(w http.ResponseWriter, event *slackevents.EventsAPIEvent, body []byte) error {
	// URL Verification
	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &r); err != nil {
			return fmt.Errorf("チャレンジレスポンスのパースに失敗: %w", err)
		}
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(r.Challenge))
		return err
	}

	// メッセージイベントの処理
	if event.Type == slackevents.CallbackEvent {
		switch ev := event.InnerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			return h.handleMessage(ev)
		}
	}

	return nil
}

func (h *Handler) handleMessage(event *slackevents.MessageEvent) error {
	switch event.Text {
	case "hello":
		_, _, err := h.client.PostMessage(
			event.Channel,
			slackapi.MsgOptionText("world", false),
		)
		if err != nil {
			return fmt.Errorf("メッセージの送信に失敗: %w", err)
		}

	case "image":
		img, err := h.imgStore.GetImage()
		if err != nil {
			return fmt.Errorf("画像の取得に失敗: %w", err)
		}

		_, err = h.client.UploadFileV2(slackapi.UploadFileV2Parameters{
			Channel:  event.Channel,
			Content:  string(img.Data),
			Filename: fmt.Sprintf("image_%d.jpg", img.ID),
			FileSize: len(img.Data),
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
