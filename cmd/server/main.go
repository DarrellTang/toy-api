package main

import (
	"log"
	"net/http"
	"os"

	"toy-api/internal/api"
)

func main() {
	addr := getAddr()

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return ":" + port
}
