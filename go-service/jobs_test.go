package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
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

// newTestJobsContext creates a JobsContext backed by an in-process miniredis instance.
func newTestJobsContext(t *testing.T) (*JobsContext, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return &JobsContext{db: &DbContext{redisClient: client}}, mr
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
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, FilterGrayscale)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestPostJob_ValidRequest_ReturnsJobID(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, FilterGrayscale)
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
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, FilterGrayscale)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}
}

func TestPostJob_JobStatusStoredInRedis(t *testing.T) {
	jc, mr := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, FilterGrayscale)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	jobID := body["id"]

	val, err := mr.Get("job:" + jobID)
	if err != nil {
		t.Fatalf("expected job:%s to exist in Redis, got error: %s", jobID, err)
	}

	var info JobInfo
	if err := json.Unmarshal([]byte(val), &info); err != nil {
		t.Fatalf("stored value is not valid JSON: %s", err)
	}

	if info.ID != jobID {
		t.Errorf("expected ID %q, got %q", jobID, info.ID)
	}
	if info.Filter != FilterGrayscale {
		t.Errorf("expected filter %q, got %q", FilterGrayscale, info.Filter)
	}
	if info.Status != "pending" {
		t.Errorf("expected status %q, got %q", "pending", info.Status)
	}
}

func TestPostJob_JobPushedToQueue(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, FilterGrayscale)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	jobID := body["id"]

	queueLen, err := jc.db.redisClient.LLen(context.Background(), "jobs:queue").Result()
	if err != nil || queueLen != 1 {
		t.Fatalf("expected 1 item in jobs:queue, got %d (err: %v)", queueLen, err)
	}

	raw, err := jc.db.redisClient.LIndex(context.Background(), "jobs:queue", 0).Result()
	if err != nil {
		t.Fatalf("failed to read queue message: %s", err)
	}

	var msg JobQueueMessage
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("queue message is not valid JSON: %s", err)
	}

	if msg.ID != jobID {
		t.Errorf("expected queue message ID %q, got %q", jobID, msg.ID)
	}
	if msg.Filter != FilterGrayscale {
		t.Errorf("expected filter %q, got %q", FilterGrayscale, msg.Filter)
	}
	if !bytes.Equal(msg.ImageData, validPNG) {
		t.Errorf("image data in queue does not match uploaded PNG")
	}
}

func TestPostJob_QueueMessageImageIsBase64Encoded(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, FilterGrayscale)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	raw, _ := jc.db.redisClient.LIndex(context.Background(), "jobs:queue", 0).Result()

	var rawMsg map[string]interface{}
	json.Unmarshal([]byte(raw), &rawMsg)

	encoded, ok := rawMsg["image_data"].(string)
	if !ok {
		t.Fatal("expected image_data to be a base64 string in JSON")
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("image_data is not valid base64: %s", err)
	}
	if !bytes.Equal(decoded, validPNG) {
		t.Error("decoded image_data does not match original PNG bytes")
	}
}

func TestPostJob_EmptyBody_Returns400(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := httptest.NewRequest("POST", "/jobs", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_MissingImage_Returns400(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, nil, FilterGrayscale)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_MissingFilter_Returns400(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, "")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_InvalidFilter_Returns400(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	req := createMultipartRequest(t, validPNG, "vaporwave")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_NotPNG_Returns400(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	notPNG := []byte("this is not a PNG file")
	req := createMultipartRequest(t, notPNG, FilterGrayscale)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostJob_UniqueJobIDs(t *testing.T) {
	jc, _ := newTestJobsContext(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", jc.handlePostJob)

	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		req := createMultipartRequest(t, validPNG, FilterGrayscale)
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
