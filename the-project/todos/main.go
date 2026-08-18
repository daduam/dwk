package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const picsumURL = "https://picsum.photos/1200"

var (
	imageClient = &http.Client{Timeout: 30 * time.Second}
	imageMu     sync.Mutex
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head><title>Todos</title></head>
<body>
<h1>Todos</h1>
<img src="/image" alt="Image expiry status" style="max-width: 400px; height: auto;">
<p>Welcome to the Todos app!</p>
</body>
</html>`))
}

func imageHandler(imageFile, imageExpiryTimestampFile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageMu.Lock()
		defer imageMu.Unlock()

		if imageExpired(imageFile, imageExpiryTimestampFile) {
			if err := fetchImage(imageFile, imageExpiryTimestampFile); err != nil {
				log.Printf("failed to fetch image: %v", err)
			}
		}

		if _, err := os.Stat(imageFile); err != nil {
			http.Error(w, "image unavailable", http.StatusServiceUnavailable)
			return
		}

		http.ServeFile(w, r, imageFile)
	}
}

func imageExpired(imageFile, imageExpiryTimestampFile string) bool {
	if _, err := os.Stat(imageFile); err != nil {
		return true
	}
	data, err := os.ReadFile(imageExpiryTimestampFile)
	if err != nil {
		return true
	}
	expiry, err := parseTimestamp(string(data))
	if err != nil {
		return true
	}
	return time.Now().After(expiry)
}

func fetchImage(imageFile, imageExpiryTimestampFile string) error {
	resp, err := imageClient.Get(picsumURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("picsum returned status %d", resp.StatusCode)
	}

	tmp, err := os.CreateTemp(filepath.Dir(imageFile), ".image-*.jpg")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), imageFile); err != nil {
		return err
	}

	expiry := time.Now().Add(10*time.Minute).Format(time.RFC3339) + "\n"
	return os.WriteFile(imageExpiryTimestampFile, []byte(expiry), 0o644)
}

func parseTimestamp(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if ts, err := time.Parse(time.RFC3339, s); err == nil {
		return ts, nil
	}
	if unix, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(unix, 0), nil
	}
	return time.Time{}, fmt.Errorf("expected RFC3339 or unix timestamp, got %q", s)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	imageExpiryTimestampFile := os.Getenv("IMAGE_EXPIRY_TIMESTAMP_FILE")
	if imageExpiryTimestampFile == "" {
		imageExpiryTimestampFile = "./image_expiry_timestamp_file.txt"
	}
	log.Println("Image expiry timestamp file: " + imageExpiryTimestampFile)

	imageFile := os.Getenv("IMAGE_FILE")
	if imageFile == "" {
		imageFile = "./image_file.jpg"
	}
	log.Println("Image file: " + imageFile)

	http.HandleFunc("GET /", indexHandler)
	http.HandleFunc("GET /image", imageHandler(imageFile, imageExpiryTimestampFile))

	log.Println("Server started in port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
