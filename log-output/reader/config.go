package main

import "os"

// Config holds the runtime settings for the reader, sourced from
// environment variables with sane defaults.
type Config struct {
	FilePath    string
	InfoPath    string
	Message     string
	PingPongURL string
	Port        string
}

// loadConfig reads configuration from environment variables, falling
// back to defaults when a variable is unset or empty.
func loadConfig() Config {
	return Config{
		FilePath:    envOr("OUTPUT_FILE", "/app/logs/output.log"),
		InfoPath:    envOr("INFORMATION_FILE", "/config/information.txt"),
		Message:     os.Getenv("MESSAGE"),
		PingPongURL: envOr("PING_PONG_URL", "http://ping-pong-svc:3456"),
		Port:        envOr("PORT", "8080"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
