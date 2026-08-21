package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var todoClient = &http.Client{Timeout: 10 * time.Second}

type Todo struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

func fetchTodos(todoAPIURL string) ([]Todo, error) {
	resp, err := todoClient.Get(todoAPIURL + "/todos")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend returned status %d", resp.StatusCode)
	}

	var todos []Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, err
	}
	return todos, nil
}
