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
)

const (
	testSigningSecret = "test-signing-secret"
)

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
