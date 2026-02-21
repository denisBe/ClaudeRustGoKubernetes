package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JobCreatedResponse struct {
	ID string `json:"id"`
}

var pngMagic = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

const (
	FilterGrayscale = "grayscale"
	FilterSepia     = "sepia"
)

var validFilters = map[string]bool{
	FilterGrayscale: true,
	FilterSepia:     false,
}

type JobsContext struct {
	db *DbContext
}

func extractImage(r *http.Request) ([]byte, string, error) {
	img, header, err := r.FormFile("image")
	if err != nil {
		return nil, "", fmt.Errorf("missing image: %w", err)
	}
	defer img.Close()

	imgBuf, err := io.ReadAll(img)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image: %w", err)
	}

	if len(imgBuf) < len(pngMagic) || !bytes.Equal(imgBuf[:len(pngMagic)], pngMagic) {
		return nil, "", fmt.Errorf("file is not a PNG image")
	}

	return imgBuf, header.Filename, nil
}

func extractFilter(r *http.Request) (string, error) {
	filter := r.FormValue("filter")
	if filter == "" {
		return "", fmt.Errorf("missing filter parameter")
	}
	if !validFilters[filter] {
		return "", fmt.Errorf("invalid filter: %s", filter)
	}
	return filter, nil
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	buf.WriteTo(w)
}

func (jc *JobsContext) handlePostJob(w http.ResponseWriter, r *http.Request) {
	imgBuf, filename, err := extractImage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter, err := extractFilter(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received file: %s (%d bytes) for filter %s\n", filename, len(imgBuf), filter)

	job, err := jc.db.CreateJob(imgBuf, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, JobCreatedResponse{ID: job.ID})
}

func (jc *JobsContext) handleGetJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET job")

	jobInfo := db.
		setHeader(w, http.StatusNotImplemented)
}
