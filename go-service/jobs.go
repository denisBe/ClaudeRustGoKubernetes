package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
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

func handlePostJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received POST job")

	img, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer img.Close()

	imgBuf, err := io.ReadAll(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(imgBuf) < len(pngMagic) || !bytes.Equal(imgBuf[:len(pngMagic)], pngMagic) {
		http.Error(w, "File is not a PNG image", http.StatusBadRequest)
		return
	}

	filter := r.FormValue("filter")
	if filter == "" {
		http.Error(w, "Missing filter parameter", http.StatusBadRequest)
		return
	}
	if !validFilters[filter] {
		http.Error(w, "Invalid filter parameter "+filter, http.StatusBadRequest)
		return
	}

	fmt.Printf("Received file: %s for filter %s\n", header.Filename, filter)

	id := uuid.New()
	var responseBody bytes.Buffer
	err = json.NewEncoder(&responseBody).Encode(JobCreatedResponse{ID: id.String()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody.Bytes())
}

func handleGetJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET job")

	setHeader(w, http.StatusNotImplemented)
}
