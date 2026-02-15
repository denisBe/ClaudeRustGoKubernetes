package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /jobs", handleGetJobs)
	mux.HandleFunc("POST /jobs", handlePostJob)
	mux.HandleFunc("GET /healthz", handleGetHealth)

	fmt.Println("Server starting on :8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

type StatusResponse struct {
	Status string `json:"status"`
}

func returnStatus(w http.ResponseWriter, status StatusResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

func handlePostJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received POST job")

	returnStatus(w, StatusResponse{Status: "Okayyy!!!"})
}

func handleGetJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET job")

	returnStatus(w, StatusResponse{Status: "Il va faire tout noir! TA GUEUUULE!!!"})
}

func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET health")

	returnStatus(w, StatusResponse{Status: "Healthy"})
}
