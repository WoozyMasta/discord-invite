// Package main implements a lightweight HTTP server that generates Discord invite links
// on demand using the Discord Bot API. The server provides logging, configurable parameters,
// and graceful shutdown handling.
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/woozymasta/discord-invite/internal/service"
)

func main() {
	if service.IsServiceMode() {
		service.RunAsService(runApp)
		return
	}

	runApp()
}

func runApp() {
	cfg := setup()

	// Setting up the HTTP server with middleware for logging
	logger := log.With().Str("component", "http").Logger()

	handler := newHandler(cfg, logger)

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       60 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       5 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info().Msgf("Starting server on %s", cfg.Listen)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Waiting for signal
	<-quit
	log.Info().Msg("Termination signal received, shutting down the server...")

	// Context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Error while shutting down the server")
	}

	log.Info().Msg("Server successfully stopped")
}
