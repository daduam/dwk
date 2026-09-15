package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

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
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server started in port " + port)
	log.Fatal(srv.ListenAndServe())
}
