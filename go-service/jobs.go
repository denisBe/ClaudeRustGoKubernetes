package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func handlePostJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received POST job")

	setHeader(w, http.StatusCreated)
	id := uuid.New()
	fmt.Printf("Generated uuid: %s", id.String())
	/*
		err := json.NewEncoder(w).Encode(status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	*/
}

func handleGetJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET job")

	setHeader(w, http.StatusNotImplemented)
}
