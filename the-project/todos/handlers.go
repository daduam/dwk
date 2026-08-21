package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
)

type indexData struct {
	Todos []Todo
	Error string
}

func indexHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := indexData{}
		todos, err := fetchTodos(cfg.TodoAPIURL)
		if err != nil {
			log.Printf("failed to fetch todos: %v", err)
			data.Error = "Failed to load todos. Please try again."
		} else {
			data.Todos = todos
		}

		var buf bytes.Buffer
		if err := indexTmpl.Execute(&buf, data); err != nil {
			log.Printf("execute index template: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		buf.WriteTo(w)
	}
}

func imageHandler(cfg config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageMu.Lock()
		defer imageMu.Unlock()

		if imageExpired(cfg.ImageFile, cfg.ImageExpiryTimestampFile) {
			if err := fetchImage(cfg.ImageFile, cfg.ImageExpiryTimestampFile); err != nil {
				log.Printf("failed to fetch image: %v", err)
			}
		}

		if _, err := os.Stat(cfg.ImageFile); err != nil {
			http.Error(w, "image unavailable", http.StatusServiceUnavailable)
			return
		}

		http.ServeFile(w, r, cfg.ImageFile)
	}
}
