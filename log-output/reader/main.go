package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	filePath := os.Getenv("OUTPUT_FILE")
	if filePath == "" {
		filePath = "/app/logs/output.log"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read output file: %v", err)
			http.Error(w, "Failed to read log file", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(data))
	})

	log.Printf("Server started at port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
