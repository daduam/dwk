package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomStr := uuid.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output := fmt.Sprintf("%s: %s", time.Now().UTC().Format(time.RFC3339), randomStr)
		fmt.Fprint(w, output)
	})

	log.Printf("Server started at port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
