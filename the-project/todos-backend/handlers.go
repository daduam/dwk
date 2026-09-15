package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type app struct {
	store *Store
}

type createTodoRequest struct {
	Content string `json:"content"`
}

func (a *app) listTodosHandler(w http.ResponseWriter, r *http.Request) {
	todos, err := a.store.List(r.Context())
	if err != nil {
		log.Printf("GET /todos: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Printf("encode todos: %v", err)
	}
}

func (a *app) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	todo, err := a.store.Add(r.Context(), req.Content)
	if err != nil {
		log.Printf("POST /todos: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		log.Printf("encode todo: %v", err)
	}
}
