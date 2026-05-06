// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kdsmith18542/pwny/internal/api"
	"github.com/kdsmith18542/pwny/internal/core"
	"github.com/kdsmith18542/pwny/internal/db"
	"github.com/kdsmith18542/pwny/internal/payload"
	"github.com/spf13/cobra"

	_ "github.com/kdsmith18542/pwny/internal/payload/encoder"
	_ "github.com/kdsmith18542/pwny/internal/payload/format"
	_ "github.com/kdsmith18542/pwny/internal/payload/stagers"
	_ "github.com/kdsmith18542/pwny/internal/payload/stages"
	_ "github.com/kdsmith18542/pwny/modules/auxiliary/scanner"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "pwny-server",
	Short: "Pwny penetration testing framework server",
	Long: `Pwny is a modern, modular penetration testing framework.
This server runs the HTTP API which clients (GUI, CLI) connect to.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file (YAML)")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("pwny-server v0.1.0")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "init-config",
		Short: "Write default config file to current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.WriteDefaultConfig("pwny.yaml")
		},
	})
}

func runServer() error {
	cfg, err := core.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	setupLogging(cfg.Log)

	database, err := db.Open(cfg.DB.Path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	eventBus := core.NewEventBus()
	sessionMgr := core.NewSessionManager()
	jobMgr := core.NewJobManager(eventBus)

	payloadGenerator := payload.NewDefaultGenerator(payload.GlobalRegistry)

	slog.Info("pwny-server starting",
		"api_addr", fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port),
		"db_path", cfg.DB.Path,
		"modules_loaded", core.ModuleCount(),
	)

	srv := api.New(api.Config{
		Host:    cfg.API.Host,
		Port:    cfg.API.Port,
		Allowed: cfg.API.Allowed,
	}, sessionMgr, jobMgr, eventBus, database, payloadGenerator)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		slog.Info("received signal, shutting down", "signal", sig)
		if err := srv.Stop(nil); err != nil {
			slog.Error("error during shutdown", "error", err)
		}
		os.Exit(0)
	}()

	return srv.Start()
}

func setupLogging(cfg core.LogConfig) {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}
