package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

func main() {
	var count atomic.Uint64

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		prev := count.Load()
		fmt.Fprintf(w, "pong %d", prev)
		count.Add(1)
	})

	log.Printf("Server started at port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
