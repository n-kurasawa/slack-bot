package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/n-kurasawa/slack-bot/internal/image"
	slackapi "github.com/slack-go/slack"
)

const (
	testSigningSecret = "test-signing-secret"
)

type mockSlackClient struct {
	messages []string
}

func (m *mockSlackClient) PostMessage(channelID string, options ...slackapi.MsgOption) (string, string, error) {
	m.messages = append(m.messages, "画像を保存しました :white_check_mark:")
	return "", "", nil
}

func (m *mockSlackClient) UploadFileV2(params slackapi.UploadFileV2Parameters) (*slackapi.FileSummary, error) {
	return nil, nil
}

type mockImageStore struct {
	savedURL string
	image    *image.Image
}

func (m *mockImageStore) GetImage() (*image.Image, error) {
	if m.image != nil {
		return m.image, nil
	}
	return &image.Image{ID: 1, URL: "https://example.com/image.jpg"}, nil
}

func (m *mockImageStore) SaveImage(url string) error {
	m.savedURL = url
	return nil
}

func TestHandler(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	tests := []struct {
		name       string
		timestamp  string
		body       []byte
		signature  string
		wantStatus int
		wantBody   string
	}{
		{
			name:      "URL検証リクエスト",
			timestamp: timestamp,
			body: func() []byte {
				event := map[string]interface{}{
					"type":      "url_verification",
					"challenge": "challenge",
				}
				b, err := json.Marshal(event)
				if err != nil {
					t.Fatal(err)
				}
				return b
			}(),
			wantStatus: http.StatusOK,
			wantBody:   "challenge",
		},
		{
			name:       "不正なJSONリクエスト",
			timestamp:  timestamp,
			body:       []byte("invalid json"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "",
		},
		{
			name:      "メッセージイベント",
			timestamp: timestamp,
			body: func() []byte {
				event := map[string]interface{}{
					"type": "event_callback",
					"event": map[string]interface{}{
						"type":    "message",
						"text":    "unknown",
						"channel": "test-channel",
					},
				}
				b, err := json.Marshal(event)
				if err != nil {
					t.Fatal(err)
				}
				return b
			}(),
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
		{
			name:       "署名が無効",
			timestamp:  timestamp,
			body:       []byte("{}"),
			signature:  "invalid",
			wantStatus: http.StatusUnauthorized,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()

			// タイムスタンプとシグネチャを設定
			req.Header.Set("X-Slack-Request-Timestamp", tt.timestamp)
			if tt.signature == "" {
				tt.signature = generateSignature(t, testSigningSecret, tt.timestamp, tt.body)
			}
			req.Header.Set("X-Slack-Signature", tt.signature)

			h := NewHandler(nil, nil, nil, testSigningSecret)
			h.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("want status %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantBody != "" && w.Body.String() != tt.wantBody {
				t.Errorf("want body %s, got %s", tt.wantBody, w.Body.String())
			}
		})
	}
}

func TestHandleMessage_UpdateImage(t *testing.T) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	mockClient := &mockSlackClient{}
	mockStore := &mockImageStore{}

	// updateImageコマンドのイベントを作成
	body, err := json.Marshal(map[string]interface{}{
		"type": "event_callback",
		"event": map[string]interface{}{
			"type":    "message",
			"text":    "updateImage https://example.com/image.jpg",
			"channel": "test-channel",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// タイムスタンプとシグネチャを設定
	req.Header.Set("X-Slack-Request-Timestamp", timestamp)
	req.Header.Set("X-Slack-Signature", generateSignature(t, testSigningSecret, timestamp, body))

	// ハンドラーを作成して実行
	h := NewHandler(mockClient, nil, mockStore, testSigningSecret)
	h.ServeHTTP(w, req)

	// レスポンスを検証
	if w.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, w.Code)
	}

	// 保存されたURLを検証
	wantURL := "https://example.com/image.jpg"
	if mockStore.savedURL != wantURL {
		t.Errorf("saved URL does not match: want %v, got %v", wantURL, mockStore.savedURL)
	}

	// 送信されたメッセージを検証
	wantMessage := "画像を保存しました :white_check_mark:"
	if len(mockClient.messages) != 1 || mockClient.messages[0] != wantMessage {
		t.Errorf("sent message does not match: want %v, got %v", wantMessage, mockClient.messages)
	}
}

func generateSignature(t *testing.T, signingSecret, timestamp string, body []byte) string {
	t.Helper()

	// バージョン番号、タイムスタンプ、ボディを結合
	baseString := fmt.Sprintf("v0:%s:%s", timestamp, body)

	// HMAC-SHA256を計算
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write([]byte(baseString))
	signature := mac.Sum(nil)

	// 16進数文字列に変換してv0:を付加
	return "v0=" + hex.EncodeToString(signature)
}
