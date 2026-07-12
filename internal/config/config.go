package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	CoolifyURL   string
	CoolifyToken string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		CoolifyURL:   strings.TrimSpace(os.Getenv("COOLIFY_URL")),
		CoolifyToken: strings.TrimSpace(os.Getenv("COOLIFY_TOKEN")),
	}

	if cfg.CoolifyURL == "" {
		return Config{}, fmt.Errorf("COOLIFY_URL is not set")
	}

	if cfg.CoolifyToken == "" {
		return Config{}, fmt.Errorf("COOLIFY_TOKEN is not set")
	}

	parsedURL, err := url.Parse(cfg.CoolifyURL)
	if err != nil {
		return Config{}, fmt.Errorf("invalid COOLIFY_URL: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return Config{}, fmt.Errorf("COOLIFY_URL must start with http:// or https://")
	}

	if parsedURL.Host == "" {
		return Config{}, fmt.Errorf("COOLIFY_URL must contain a valid host")
	}

	cfg.CoolifyURL = strings.TrimRight(cfg.CoolifyURL, "/")

	return cfg, nil
}
