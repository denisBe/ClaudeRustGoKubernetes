package main

import (
	"net/http"
)

func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	setHeader(w, http.StatusOK)
}
