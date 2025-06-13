package main

import (
	"frames/imageprocessing"
	"log"
	"net/http"
)

func main() {
	// Log when the server starts
	log.Println("Starting server on port :8080...")

	// Set up the overlay handler
	http.HandleFunc("/imageprocessing/overlay", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:", r.URL.String())
		imageprocessing.OverlayHandler(w, r)
		log.Println("Request processed successfully")
	})

	// Start the server and log if it fails to start
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
