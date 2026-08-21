package main

import (
	"html/template"
	"log"
	"net/http"
)

var indexTmpl *template.Template

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("config: port=%s image=%s expiry=%s todo-api=%s",
		cfg.Port, cfg.ImageFile, cfg.ImageExpiryTimestampFile, cfg.TodoAPIURL)

	indexTmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	http.HandleFunc("GET /", indexHandler(cfg))
	http.HandleFunc("GET /image", imageHandler(cfg))

	log.Println("Server started in port " + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
