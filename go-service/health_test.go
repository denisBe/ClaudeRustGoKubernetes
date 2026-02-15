package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint_ReturnsOK(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleGetHealth)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestHealthEndpoint_ReturnsJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleGetHealth)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}
}

func TestHealthEndpoint_ReturnsHealthyStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleGetHealth)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	var body StatusResponse
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response body: %s", err)
	}

	if body.Status == "" {
		t.Error("expected non-empty status field")
	}
}

func TestHealthEndpoint_WrongMethodReturns405(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleGetHealth)

	req := httptest.NewRequest("POST", "/healthz", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405 for POST /healthz, got %d", rec.Code)
	}
}
