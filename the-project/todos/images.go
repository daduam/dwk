package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	imageClient = &http.Client{Timeout: 30 * time.Second}
	imageMu     sync.Mutex
)

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

func fetchImage(picsumURL, imageFile, imageExpiryTimestampFile string) error {
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
