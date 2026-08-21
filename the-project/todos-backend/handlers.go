package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type createTodoRequest struct {
	Content string `json:"content"`
}

func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(store.List()); err != nil {
		log.Printf("encode todos: %v", err)
	}
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	todo := store.Add(req.Content)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("encode todo: %v", err)
	}
}
