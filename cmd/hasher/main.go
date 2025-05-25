// Package main provides entry-point for application. Determines application
// mode, setups logger, initializes configuration and bootstrapps application.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmybsv/leadgen-test-task/internal/app"
	"github.com/tmybsv/leadgen-test-task/internal/infrastructure/config"
)

func main() {
	mode := config.Mode(os.Getenv("HASHER_MODE"))
	if mode == "" {
		mode = config.ModeDevelopment
	}

	log := setupLogger(mode)

	cfg, err := config.New(mode, log)
	if err != nil {
		log.Error("failed to create new config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.New(cfg, log)
	defer app.Stop()

	errCh := make(chan error, 1)
	go func() {
		if err := app.GRPCServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to run gRPC server: %w", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sigCh:
		log.Info("shutdown signal received", slog.String("signal", s.String()))
		cancel()
	case err := <-errCh:
		log.Error("failed to bootstrap application", slog.String("error", err.Error()))
		cancel()
	case <-ctx.Done():
	}
}

func setupLogger(mode config.Mode) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	var handler slog.Handler
	if mode == config.ModeProduction {
		opts.Level = slog.LevelInfo
		opts.AddSource = false
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
