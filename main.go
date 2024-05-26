package main

import (
	"log"
	"net/http"
	"os"
	"sykell-back/api"
)

func main() {
	// create an API
	r := api.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
		log.Println("PORT environment variable not set, defaulting to 8080")
	} else {
		log.Println("PORT environment variable found, using port " + port)
	}

	// start server & check errors
	address := ":" + port
	log.Println("Server is running on port " + port)
	err := http.ListenAndServe(address, api.HandleCORS(r))
	if err != nil {
		if os.IsPermission(err) {
			log.Fatalf("Permission error while trying to listen on port %s: %v", port, err)
		} else if os.IsTimeout(err) {
			log.Fatalf("Timeout error while trying to listen on port %s: %v", port, err)
		} else {
			log.Fatalf("Error starting server on port %s: %v", port, err)
		}
	}
}
