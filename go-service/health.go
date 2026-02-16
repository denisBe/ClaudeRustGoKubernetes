package main

import (
	"fmt"
	"net/http"
)

func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET health")

	setHeader(w, http.StatusOK)
}
