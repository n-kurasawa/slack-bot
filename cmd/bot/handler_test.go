package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLVerification(t *testing.T) {
	challenge := "test_challenge_token"
	event := SlackEvent{
		Token:     "test_token",
		Challenge: challenge,
		Type:      "url_verification",
	}

	body, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("JSONのエンコードに失敗: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/slack/events", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler := createHandler(nil, nil) // テスト用にクライアントとDBをnilに設定
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期待したステータスコード %d, 実際のステータスコード %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != challenge {
		t.Errorf("期待したレスポンス %s, 実際のレスポンス %s", challenge, w.Body.String())
	}
}

func TestInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/slack/events", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler := createHandler(nil, nil)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("期待したステータスコード %d, 実際のステータスコード %d", http.StatusBadRequest, w.Code)
	}
}
