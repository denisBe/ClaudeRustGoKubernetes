package main

import (
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func registerRoutes(mux *http.ServeMux, redisClient *redis.Client) {
	jobsContext := &JobsContext{redisClient: redisClient}
	mux.HandleFunc("GET /healthz", handleGetHealth)
	mux.HandleFunc("GET /jobs", jobsContext.handleGetJobs)
	mux.HandleFunc("POST /jobs", jobsContext.handlePostJob)

	staticFS, err := StaticFS()
	if err != nil {
		log.Fatal("Failed to load embedded static files:", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	// Serve static assets (CSS, JS) at /static/
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Serve index.html at root
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFileFS(w, r, staticFS, "index.html")
	})
}
