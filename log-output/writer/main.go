package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomStr := uuid.New().String()

	filePath := os.Getenv("OUTPUT_FILE")
	if filePath == "" {
		filePath = "/app/logs/output.log"
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer f.Close()

	log.Printf("Writing to %s every 5 seconds", filePath)

	for {
		msg := fmt.Sprintf("%s: %s", time.Now().UTC().Format(time.RFC3339), randomStr)
		log.Print(msg)
		if _, err := fmt.Fprintln(f, msg); err != nil {
			log.Fatalf("Failed to write: %v", err)
		}
		time.Sleep(5 * time.Second)
	}
}
