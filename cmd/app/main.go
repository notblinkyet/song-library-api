package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/notblinkyet/song-library-api/internal/config"
	"github.com/notblinkyet/song-library-api/internal/database/postgresql"
	"github.com/notblinkyet/song-library-api/internal/lib/api"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
	"github.com/notblinkyet/song-library-api/internal/logger"
	"github.com/notblinkyet/song-library-api/internal/services"
	myHttp "github.com/notblinkyet/song-library-api/internal/transport/http"
)

// @title Song Library API
// @version 1.0
// @description A simple RESTful API for managing a song library.

// @host localhost:9090
// @BasePath /

func main() {
	// Load configuration
	config := config.MustLoadConfig()
	log := logger.SetupLogger()

	log.Info("Configuration loaded", slog.Any("config", config))

	// Set up signal handling for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Initialize dependencies
	log.Info("Initializing dependencies")
	db, err := postgresql.NewPostgreSQL(config)
	if err != nil {
		log.Error("Failed to connect to database", sl.Error(err))
		os.Exit(1)
	}
	log.Info("Database connection established")

	apiClient := api.NewApiClient(config.ApiAddrURL)
	server := services.NewSongLibraryService(db, apiClient, log)
	handler := myHttp.NewHandler(server, log)

	// Set up HTTP router and endpoints
	r := chi.NewMux()
	handler.FillEndpoints(r)
	log.Debug("HTTP routes configured")

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort),
		Handler:      r,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		IdleTimeout:  config.IdleTimeout,
	}
	log.Info("HTTP server configured", slog.String("address", srv.Addr))

	// Run server in a separate goroutine
	go func() {
		log.Info("Starting server", slog.String("address", srv.Addr))
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", sl.Error(err))
			done <- syscall.SIGTERM
		}
	}()
	log.Info("Server is running")

	// Wait for shutdown signal
	<-done
	log.Info("Shutdown signal received, stopping server")

	// Gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Error("Failed to shut down server", sl.Error(err))
	} else {
		log.Info("Server stopped gracefully")
	}

	// Close database connection
	db.Close()
	log.Info("Database connection closed")
}
