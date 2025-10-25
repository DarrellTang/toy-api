package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	healthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected content-type application/json, got %q", ct)
	}

	var payload healthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}

	if payload.Status != "ok" {
		t.Fatalf("expected status 'ok', got %q", payload.Status)
	}
}

func TestMessageHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/message", nil)
	rr := httptest.NewRecorder()

	messageHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var payload messageResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected error decoding response: %v", err)
	}

	expected := "Welcome to the toy API"
	if payload.Message != expected {
		t.Fatalf("expected message %q, got %q", expected, payload.Message)
	}
}
