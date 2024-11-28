package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/notblinkyet/song-library-api/internal/config"
	"github.com/notblinkyet/song-library-api/internal/lib/sl"
	"github.com/notblinkyet/song-library-api/internal/logger"
)

func main() {
	// Define rollback flag for migrations
	var rollback int
	flag.IntVar(&rollback, "rollback", 0, "number of steps to rollback")
	flag.Parse()

	// Load configuration
	cfg := config.MustLoadConfig()

	// Set up logging
	log := logger.SetupLogger()
	log.Info("Loaded configuration", slog.Any("config", cfg))
	log.Info("Rollback steps requested", slog.Int("rollback", rollback))

	// Prepare database connection string and migration path
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)
	migrationPath := cfg.MigrationPath

	if connString == "" {
		log.Error("Database connection string is empty")
		panic("Database connection string is required")
	}

	if migrationPath == "" {
		log.Error("Migration path is empty")
		panic("Migration path is required")
	}

	// Initialize migration engine
	m, err := migrate.New(
		"file://"+migrationPath,
		fmt.Sprintf("%s?x-migrations-table=", connString),
	)
	if err != nil {
		log.Error("Failed to create migration engine", sl.Error(err))
		panic(err)
	}
	log.Info("Migration engine initialized successfully")

	// Handle rollback logic
	if rollback > 0 {
		log.Info("Rolling back migrations", slog.Int("steps", rollback))
		if err = m.Steps(-rollback); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info("No migrations to rollback")
				return
			}
			log.Error("Failed to rollback migrations", sl.Error(err))
			panic(err)
		}
		log.Info("Rollback completed successfully", slog.Int("steps", rollback))
		return
	}

	// Apply migrations
	log.Info("Applying migrations")
	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No migrations to apply")
			return
		}
		log.Error("Failed to apply migrations", sl.Error(err))
		panic(err)
	}
	log.Info("Migrations applied successfully")
}
