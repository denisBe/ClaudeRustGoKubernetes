package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Valid PNG: smallest possible 1x1 pixel PNG file
var validPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, // PNG header
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 pixel
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
	0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41,
	0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
	0x00, 0x00, 0x02, 0x00, 0x01, 0xe2, 0x21, 0xbc,
	0x33, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
	0x44, 0xae, 0x42, 0x60, 0x82,
}

// createMultipartRequest builds a multipart/form-data request with an image part and a filter part.
func createMultipartRequest(t *testing.T, png []byte, filter string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if png != nil {
		part, err := w.CreateFormFile("image", "test.png")
		if err != nil {
			t.Fatalf("failed to create image part: %s", err)
		}
		_, err = io.Copy(part, bytes.NewReader(png))
		if err != nil {
			t.Fatalf("failed to write image part: %s", err)
		}
	}

	if filter != "" {
		err := w.WriteField("filter", filter)
		if err != nil {
			t.Fatalf("failed to write filter field: %s", err)
		}
	}

	err := w.Close()
	if err != nil {
		t.Fatalf("failed to close multipart writer: %s", err)
	}

	req := httptest.NewRequest("POST", "/jobs", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func TestPostJob_ValidRequest_Returns201(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	req := createMultipartRequest(t, validPNG, "grayscale")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestPostJob_ValidRequest_ReturnsJobID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	req := createMultipartRequest(t, validPNG, "sepia")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	var body map[string]string
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("response is not valid JSON: %s", err)
	}

	id, exists := body["id"]
	if !exists || id == "" {
		t.Error("expected non-empty 'id' field in response")
	}
}

func TestPostJob_ValidRequest_ReturnsJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	req := createMultipartRequest(t, validPNG, "grayscale")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}
}

func TestPostJob_EmptyBody_Returns400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	req := httptest.NewRequest("POST", "/jobs", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_MissingImage_Returns400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	req := createMultipartRequest(t, nil, "grayscale")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_MissingFilter_Returns400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	req := createMultipartRequest(t, validPNG, "")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_InvalidFilter_Returns400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	req := createMultipartRequest(t, validPNG, "vaporwave")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_NotPNG_Returns400(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	notPNG := []byte("this is not a PNG file")
	req := createMultipartRequest(t, notPNG, "grayscale")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_UniqueJobIDs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", handlePostJob)

	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		req := createMultipartRequest(t, validPNG, "grayscale")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		var body map[string]string
		json.NewDecoder(rec.Body).Decode(&body)

		id := body["id"]
		if ids[id] {
			t.Errorf("duplicate job ID: %s", id)
		}
		ids[id] = true
	}
}
