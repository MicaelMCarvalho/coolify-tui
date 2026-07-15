package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/term"
)

const configFileName = "config.json"

func Path() (string, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(
			"find user config directory: %w",
			err,
		)
	}

	return filepath.Join(
		configDirectory,
		"coolify-tui",
		configFileName,
	), nil
}

func loadSavedConfig() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}

	if err != nil {
		return Config{}, fmt.Errorf(
			"read configuration file: %w",
			err,
		)
	}

	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf(
			"decode configuration file %s: %w",
			path,
			err,
		)
	}

	return cfg, nil
}

func Save(cfg Config) (string, error) {
	validated, err := Validate(cfg)
	if err != nil {
		return "", err
	}

	path, err := Path()
	if err != nil {
		return "", err
	}

	directory := filepath.Dir(path)

	if err := os.MkdirAll(directory, 0700); err != nil {
		return "", fmt.Errorf(
			"create configuration directory: %w",
			err,
		)
	}

	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0600,
	)
	if err != nil {
		return "", fmt.Errorf(
			"open configuration file: %w",
			err,
		)
	}

	// Correct permissions if this file already existed.
	if err := file.Chmod(0600); err != nil {
		_ = file.Close()

		return "", fmt.Errorf(
			"protect configuration file: %w",
			err,
		)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(validated); err != nil {
		_ = file.Close()

		return "", fmt.Errorf(
			"write configuration file: %w",
			err,
		)
	}

	if err := file.Close(); err != nil {
		return "", fmt.Errorf(
			"close configuration file: %w",
			err,
		)
	}

	return path, nil
}

func Configure() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Coolify URL: ")

	coolifyURL, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("read Coolify URL: %w", err)
	}

	fmt.Print("Coolify API token: ")

	tokenBytes, err := term.ReadPassword(os.Stdin.Fd())

	// Move to a new line after the hidden token prompt.
	fmt.Println()

	if err != nil {
		return fmt.Errorf("read Coolify API token: %w", err)
	}

	cfg := Config{
		CoolifyURL:   strings.TrimSpace(coolifyURL),
		CoolifyToken: strings.TrimSpace(string(tokenBytes)),
	}

	path, err := Save(cfg)
	if err != nil {
		return err
	}

	fmt.Printf("Configuration saved to %s\n", path)

	return nil
}
