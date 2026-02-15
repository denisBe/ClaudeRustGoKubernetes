package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	registerRoutes(mux)

	fmt.Println("Server starting on :8081")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
