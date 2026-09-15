package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type app struct {
	store *Store
}

func (a *app) rootHandler(w http.ResponseWriter, r *http.Request) {
	n, err := a.store.Increment(r.Context())
	if err != nil {
		log.Printf("GET /: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "pong %d\n", n)
}

func (a *app) pingsHandler(w http.ResponseWriter, r *http.Request) {
	n, err := a.store.Count(r.Context())
	if err != nil {
		log.Printf("GET /pings: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d\n", n)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	store, err := NewStore(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer store.Close()

	a := &app{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.rootHandler)
	mux.HandleFunc("GET /pings", a.pingsHandler)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server started at port %s", port)
	log.Fatal(srv.ListenAndServe())
}
