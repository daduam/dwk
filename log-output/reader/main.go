package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func fetchCount(client *http.Client, baseURL string) string {
	resp, err := client.Get(strings.TrimRight(baseURL, "/") + "/pings")
	if err != nil {
		log.Printf("Failed to fetch requests count from %s: %v", baseURL, err)
		return "0"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body from %s: %v", baseURL, err)
		return "0"
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ping-pong app returned status %d: %s", resp.StatusCode, body)
		return "0"
	}

	n, err := strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
	if err != nil {
		log.Printf("Invalid requests count from %s: %q", baseURL, body)
		return "0"
	}
	return strconv.FormatUint(n, 10)
}

func logHandler(filePath string, client *http.Client, pingPongURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read output file: %v", err)
			http.Error(w, "Failed to read log file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s\nPing / Pongs: %s\n", string(data), fetchCount(client, pingPongURL))
	}
}

func main() {
	filePath := os.Getenv("OUTPUT_FILE")
	if filePath == "" {
		filePath = "/app/logs/output.log"
	}

	pingPongURL := os.Getenv("PING_PONG_URL")
	if pingPongURL == "" {
		pingPongURL = "http://ping-pong-svc:3456"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{Timeout: 5 * time.Second}

	http.HandleFunc("GET /{$}", logHandler(filePath, client, pingPongURL))

	log.Printf("Server started at port %s, ping-pong base URL %s", port, pingPongURL)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
