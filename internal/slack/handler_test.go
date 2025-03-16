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
	savedName string
	savedURL  string
	image     *image.Image
}

func (m *mockImageStore) GetImage() (*image.Image, error) {
	if m.image != nil {
		return m.image, nil
	}
	return &image.Image{ID: 1, URL: "https://example.com/image.jpg", Name: "test"}, nil
}

func (m *mockImageStore) GetImageByName(name string) (*image.Image, error) {
	if m.image != nil {
		return m.image, nil
	}
	return &image.Image{ID: 1, URL: "https://example.com/image.jpg", Name: name}, nil
}

func (m *mockImageStore) SaveImage(name, url string) error {
	m.savedName = name
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
	tests := []struct {
		name       string
		text       string
		wantName   string
		wantURL    string
		wantStatus int
		wantError  bool
	}{
		{
			name:       "画像の保存（正常系）",
			text:       "updateImage cat https://example.com/image.jpg",
			wantName:   "cat",
			wantURL:    "https://example.com/image.jpg",
			wantStatus: http.StatusOK,
			wantError:  false,
		},
		{
			name:       "画像の保存（引数不足）",
			text:       "updateImage cat",
			wantName:   "",
			wantURL:    "",
			wantStatus: http.StatusOK,
			wantError:  true,
		},
		{
			name:       "画像の保存（引数なし）",
			text:       "updateImage",
			wantName:   "",
			wantURL:    "",
			wantStatus: http.StatusOK,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			mockClient := &mockSlackClient{}
			mockStore := &mockImageStore{}

			// イベントを作成
			body, err := json.Marshal(map[string]interface{}{
				"type": "event_callback",
				"event": map[string]interface{}{
					"type":    "message",
					"text":    tt.text,
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
			if w.Code != tt.wantStatus {
				t.Errorf("want status %d, got %d", tt.wantStatus, w.Code)
			}

			// 保存された値を検証
			if tt.wantName != "" && mockStore.savedName != tt.wantName {
				t.Errorf("saved name does not match: want %v, got %v", tt.wantName, mockStore.savedName)
			}
			if tt.wantURL != "" && mockStore.savedURL != tt.wantURL {
				t.Errorf("saved URL does not match: want %v, got %v", tt.wantURL, mockStore.savedURL)
			}

			// エラーケースの検証
			if tt.wantError {
				if len(mockClient.messages) > 0 {
					t.Error("エラー時にメッセージが送信されてはいけません")
				}
			} else {
				// 成功時のメッセージを検証
				wantMessage := "画像を保存しました :white_check_mark:"
				if len(mockClient.messages) != 1 || mockClient.messages[0] != wantMessage {
					t.Errorf("sent message does not match: want %v, got %v", wantMessage, mockClient.messages)
				}
			}
		})
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
