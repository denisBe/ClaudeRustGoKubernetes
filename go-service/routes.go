package main

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", handleGetHealth)
	mux.HandleFunc("GET /jobs", handleGetJobs)
	mux.HandleFunc("POST /jobs", handlePostJob)
}
