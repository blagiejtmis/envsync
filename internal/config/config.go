// Package config handles loading and validating envsync configuration
// from a TOML file (typically .envsync.toml in the project root).
package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds the top-level envsync configuration.
type Config struct {
	// Store is the backend store configuration.
	Store StoreConfig `toml:"store"`
	// Env is the environment file configuration.
	Env EnvConfig `toml:"env"`
}

// StoreConfig describes where the encrypted secret store lives.
type StoreConfig struct {
	// Path is the file system path to the store directory or file.
	Path string `toml:"path"`
}

// EnvConfig describes which .env file envsync manages.
type EnvConfig struct {
	// File is the path to the .env file (default: ".env").
	File string `toml:"file"`
}

const DefaultConfigFile = ".envsync.toml"

// Load reads the config file at path and returns a validated Config.
// If path is empty, DefaultConfigFile is used.
func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigFile
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks that required fields are present.
func (c *Config) validate() error {
	if c.Store.Path == "" {
		return errors.New("config: store.path must not be empty")
	}
	if c.Env.File == "" {
		c.Env.File = ".env"
	}
	return nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Store: StoreConfig{Path: ".envsync-store"},
		Env:   EnvConfig{File: ".env"},
	}
}
