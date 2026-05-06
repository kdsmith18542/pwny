// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	API    APIConfig    `mapstructure:"api"`
	DB     DBConfig     `mapstructure:"database"`
	Log    LogConfig    `mapstructure:"logging"`
	Module ModuleConfig `mapstructure:"module"`
}

type APIConfig struct {
	Host    string   `mapstructure:"host"`
	Port    int      `mapstructure:"port"`
	Allowed []string `mapstructure:"allowed_origins"`
	TLSKey  string   `mapstructure:"tls_key"`
	TLSCert string   `mapstructure:"tls_cert"`
}

type DBConfig struct {
	Path    string `mapstructure:"path"`
	CredKey string `mapstructure:"credential_encryption_key"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type ModuleConfig struct {
	Paths []string `mapstructure:"paths"`
}

func DefaultConfig() Config {
	return Config{
		API: APIConfig{
			Host:    "127.0.0.1",
			Port:    31337,
			Allowed: []string{"http://localhost:1420", "tauri://localhost"},
		},
		DB: DBConfig{
			Path: filepath.Join(os.TempDir(), "pwny", "pwny.db"),
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("pwny")
	v.SetConfigType("yaml")

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.pwny")
		v.AddConfigPath("/etc/pwny")
	}

	v.SetEnvPrefix("PWNY")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
		slog.Info("no config file found, using defaults", "path", path)
	}

	cfg := DefaultConfig()
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	dir := filepath.Dir(cfg.DB.Path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("error creating db directory: %w", err)
	}

	return &cfg, nil
}

func WriteDefaultConfig(path string) error {
	cfg := DefaultConfig()
	v := viper.New()
	v.Set("api", cfg.API)
	v.Set("database", cfg.DB)
	v.Set("logging", cfg.Log)
	v.Set("module", cfg.Module)
	return v.WriteConfigAs(path)
}
