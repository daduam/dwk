package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started in port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
