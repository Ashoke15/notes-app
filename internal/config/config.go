package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	VaultDir string `json:"vault_dir,omitempty"`
	Theme    string `json:"theme,omitempty"`
}

func Default() Config {
	return Config{Theme: "auto"}
}

func Puth() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("Get ConfigDir: %w",err)
	}

	return filepath.Join(dir, "totion", "config.json"), nil
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Default(), nil
	}

	if err != nil {
		return Config{}, fmt.Errorf("Read Config err: %w", err)
	}

	var cfg Config
	if err :=json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("Parse Comfig err: %w",err)
	}

	if cfg.Theme == ""{
		cfg.Theme = "auto"
	}

	return cfg, nil	
}

func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("creat config directory: %w",err)
	}

	data, err := json.MarshalIndent(cfg, ""," ")
	if err != nil {
		return fmt.Errorf("encode config: %w",err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config: %w",err)
	}

	return nil
}

func ExpandHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}

	return filepath.Join(home, strings.TrimPrefix(path, "~")), nil
}