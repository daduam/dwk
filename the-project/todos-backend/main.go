package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func withRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("request: %s %s (%s)", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
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
	mux.HandleFunc("GET /todos", a.listTodosHandler)
	mux.HandleFunc("POST /todos", a.createTodoHandler)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           withRequestLogging(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server started in port " + port)
	log.Fatal(srv.ListenAndServe())
}
