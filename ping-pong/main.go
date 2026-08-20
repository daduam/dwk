package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func rootHandler(mu *sync.Mutex, countFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		var n uint64
		if data, err := os.ReadFile(countFile); err == nil {
			if parsed, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64); err == nil {
				n = parsed
			}
		}
		n++
		if err := os.WriteFile(countFile, []byte(strconv.FormatUint(n, 10)), 0o644); err != nil {
			log.Printf("failed to write request count to %s: %v", countFile, err)
		}
		fmt.Fprintf(w, "pong %d\n", n)
	}
}

func pingsHandler(mu *sync.Mutex, countFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		var n uint64
		if data, err := os.ReadFile(countFile); err == nil {
			if parsed, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64); err == nil {
				n = parsed
			}
		}
		fmt.Fprintf(w, "%d\n", n)
	}
}

func main() {
	var mu sync.Mutex

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	countFile := os.Getenv("REQUESTS_COUNT_FILE")
	if countFile == "" {
		countFile = "requests_count.txt"
	}

	http.HandleFunc("GET /{$}", rootHandler(&mu, countFile))
	http.HandleFunc("GET /pings", pingsHandler(&mu, countFile))

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
