package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		wantStatus int
		wantBody   string
	}{
		{
			name: "URL検証リクエスト",
			body: func() []byte {
				event := Event{
					Token:     "token",
					Challenge: "challenge",
					Type:      "url_verification",
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
			body:       []byte("invalid json"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "",
		},
		{
			name: "メッセージイベント",
			body: func() []byte {
				event := Event{
					Type: "event_callback",
					Event: struct {
						Type    string `json:"type"`
						Text    string `json:"text"`
						Channel string `json:"channel"`
						User    string `json:"user"`
					}{
						Type:    "message",
						Text:    "unknown",
						Channel: "test-channel",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()

			h := NewHandler(nil, nil, nil)
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
