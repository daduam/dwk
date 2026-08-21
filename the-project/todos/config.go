package main

import (
	"errors"
	"os"
)

type config struct {
	Port                     string
	ImageFile                string
	ImageExpiryTimestampFile string
	TodoAPIURL               string
}

func loadConfig() (config, error) {
	cfg := config{
		Port:                     envOrDefault("PORT", "8080"),
		ImageFile:                envOrDefault("IMAGE_FILE", "./image_file.jpg"),
		ImageExpiryTimestampFile: envOrDefault("IMAGE_EXPIRY_TIMESTAMP_FILE", "./image_expiry_timestamp_file.txt"),
		TodoAPIURL:               os.Getenv("TODO_API_URL"),
	}
	if cfg.TodoAPIURL == "" {
		return config{}, errors.New("TODO_API_URL is required")
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
