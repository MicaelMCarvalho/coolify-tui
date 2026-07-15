package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	CoolifyURL   string `json:"coolify_url"`
	CoolifyToken string `json:"coolify_token"`
}

func Load() (Config, error) {
	// Load .env values only when the variables are not already exported.
	_ = godotenv.Load()

	cfg := Config{
		CoolifyURL:   strings.TrimSpace(os.Getenv("COOLIFY_URL")),
		CoolifyToken: strings.TrimSpace(os.Getenv("COOLIFY_TOKEN")),
	}

	// Use the saved configuration only for missing values.
	if cfg.CoolifyURL == "" || cfg.CoolifyToken == "" {
		saved, err := loadSavedConfig()
		if err != nil {
			return Config{}, err
		}

		if cfg.CoolifyURL == "" {
			cfg.CoolifyURL = saved.CoolifyURL
		}

		if cfg.CoolifyToken == "" {
			cfg.CoolifyToken = saved.CoolifyToken
		}
	}

	return Validate(cfg)
}

func Validate(cfg Config) (Config, error) {
	cfg.CoolifyURL = strings.TrimSpace(cfg.CoolifyURL)
	cfg.CoolifyToken = strings.TrimSpace(cfg.CoolifyToken)

	if cfg.CoolifyURL == "" {
		return Config{}, fmt.Errorf(
			"COOLIFY_URL is not set; run `coolify-tui configure`",
		)
	}

	if cfg.CoolifyToken == "" {
		return Config{}, fmt.Errorf(
			"COOLIFY_TOKEN is not set; run `coolify-tui configure`",
		)
	}

	parsedURL, err := url.Parse(cfg.CoolifyURL)
	if err != nil {
		return Config{}, fmt.Errorf("invalid COOLIFY_URL: %w", err)
	}

	if parsedURL.Scheme != "http" &&
		parsedURL.Scheme != "https" {
		return Config{}, fmt.Errorf(
			"COOLIFY_URL must start with http:// or https://",
		)
	}

	if parsedURL.Host == "" {
		return Config{}, fmt.Errorf(
			"COOLIFY_URL must contain a valid host",
		)
	}

	cfg.CoolifyURL = strings.TrimRight(cfg.CoolifyURL, "/")

	return cfg, nil
}
