package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestReturnStatus_SetsContentTypeJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	returnStatus(rec, StatusResponse{Status: "ok"})

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}
}

func TestReturnStatus_SetsStatusCode200(t *testing.T) {
	rec := httptest.NewRecorder()

	returnStatus(rec, StatusResponse{Status: "ok"})

	if rec.Code != 200 {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestReturnStatus_WritesValidJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	returnStatus(rec, StatusResponse{Status: "test-value"})

	var body StatusResponse
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("response is not valid JSON: %s", err)
	}

	if body.Status != "test-value" {
		t.Errorf("expected status %q, got %q", "test-value", body.Status)
	}
}
