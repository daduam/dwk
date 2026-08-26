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

func logHandler(cfg Config, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := os.ReadFile(cfg.InfoPath)
		if err != nil {
			log.Printf("Failed to read information file: %v", err)
			http.Error(w, "Failed to read information file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "file content: %s\nenv variable: MESSAGE=%s\n", info, cfg.Message)

		data, err := os.ReadFile(cfg.FilePath)
		if err != nil {
			log.Printf("Failed to read output file: %v", err)
			http.Error(w, "Failed to read log file", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s\nPing / Pongs: %s\n", string(data), fetchCount(client, cfg.PingPongURL))
	}
}

func main() {
	cfg := loadConfig()

	client := &http.Client{Timeout: 5 * time.Second}

	http.HandleFunc("GET /{$}", logHandler(cfg, client))

	log.Printf("Server started at port %s, ping-pong base URL %s", cfg.Port, cfg.PingPongURL)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
