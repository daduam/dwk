package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func readCount(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read requests count file %s: %v", path, err)
		return "0"
	}
	n, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		log.Printf("Invalid requests count in %s: %q", path, data)
		return "0"
	}
	return strconv.FormatUint(n, 10)
}

func logHandler(filePath, countFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read output file: %v", err)
			http.Error(w, "Failed to read log file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s\nPing / Pongs: %s\n", string(data), readCount(countFile))
	}
}

func main() {
	filePath := os.Getenv("OUTPUT_FILE")
	if filePath == "" {
		filePath = "/app/logs/output.log"
	}

	countFile := os.Getenv("REQUESTS_COUNT_FILE")
	if countFile == "" {
		countFile = "requests_count.txt"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("GET /", logHandler(filePath, countFile))

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
