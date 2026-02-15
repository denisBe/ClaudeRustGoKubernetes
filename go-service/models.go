package main

import (
	"encoding/json"
	"net/http"
)

type StatusResponse struct {
	Status string `json:"status"`
}

func returnStatus(w http.ResponseWriter, status StatusResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
