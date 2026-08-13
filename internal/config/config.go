package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Addr, PluginDir, DataDir string
	LogLevel                 slog.Level
	PluginTimeout            time.Duration
}

func Load() (Config, error) {
	dataDir := value("CMDRY_DATA_DIR", "/opt/cmdry/data")
	cfg := Config{Addr: value("CMDRY_ADDR", "127.0.0.1:8080"), PluginDir: value("CMDRY_PLUGIN_DIR", "/opt/cmdry/plugins"), DataDir: dataDir, PluginTimeout: 8 * time.Second}
	switch strings.ToLower(value("CMDRY_LOG_LEVEL", "info")) {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
	case "info":
		cfg.LogLevel = slog.LevelInfo
	case "warn", "warning":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	default:
		return Config{}, fmt.Errorf("invalid CMDRY_LOG_LEVEL")
	}
	if err := os.MkdirAll(filepath.Clean(dataDir), 0750); err != nil {
		return Config{}, fmt.Errorf("create data directory: %w", err)
	}
	return cfg, nil
}
func value(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
