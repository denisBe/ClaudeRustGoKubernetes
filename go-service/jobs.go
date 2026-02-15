package main

import (
	"fmt"
	"net/http"
)

func handlePostJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received POST job")

	returnStatus(w, StatusResponse{Status: "Okayyy!!!"})
}

func handleGetJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET job")

	returnStatus(w, StatusResponse{Status: "Il va faire tout noir! TA GUEUUULE!!!"})
}
