package database

import (
	"embed"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// RunMigrations runs database migrations using the provided DSN and embedded filesystem.
// It automatically converts postgres:// or postgresql:// schemes to pgx5:// for the migrate tool.
func RunMigrations(dsn string, fs embed.FS, path string) error {
	migrationDSN := dsn
	if strings.HasPrefix(dsn, "postgres://") {
		migrationDSN = strings.Replace(dsn, "postgres://", "pgx5://", 1)
	} else if strings.HasPrefix(dsn, "postgresql://") {
		migrationDSN = strings.Replace(dsn, "postgresql://", "pgx5://", 1)
	}

	d, err := iofs.New(fs, path)
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, migrationDSN)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		slog.Info("no new migrations to apply")
	} else {
		slog.Info("migrations completed successfully")
	}

	return nil
}
