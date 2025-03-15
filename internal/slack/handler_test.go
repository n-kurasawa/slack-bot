package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLVerification(t *testing.T) {
	event := Event{
		Token:     "token",
		Challenge: "challenge",
		Type:      "url_verification",
	}

	body, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h := NewHandler(nil, nil, nil)
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != event.Challenge {
		t.Errorf("want %s, got %s", event.Challenge, w.Body.String())
	}
}

func TestInvalidJSON(t *testing.T) {
	body := bytes.NewReader([]byte("invalid json"))
	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()

	h := NewHandler(nil, nil, nil)
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusBadRequest, w.Code)
	}
}
