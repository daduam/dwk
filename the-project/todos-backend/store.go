package main

import (
	"slices"
	"sync"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

type todoStore struct {
	mu     sync.Mutex
	todos  []Todo
	nextID int
}

var store = &todoStore{nextID: 1}

func init() {
	store.Add("Learn Kubernetes")
	store.Add("Deploy the todos app")
	store.Add("Write a README")
}

func (s *todoStore) List() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Todo, len(s.todos))
	copy(out, s.todos)
	slices.SortFunc(out, func(a, b Todo) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})
	return out
}

func (s *todoStore) Add(content string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := Todo{
		ID:        s.nextID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.nextID++
	s.todos = append(s.todos, t)
	return t
}
