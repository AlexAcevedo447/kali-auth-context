package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"kali-auth-context/internal/infrastructure/config"
)

func NewPool(cfg *config.Config) (*bun.DB, error) {
	if err := ensureDatabaseExists(cfg); err != nil {
		return nil, fmt.Errorf("ensure database exists: %w", err)
	}

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(databaseUrl),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	return db, nil
}

// ensureDatabaseExists connects to the "postgres" maintenance database and creates
// the target database if it does not already exist.
func ensureDatabaseExists(cfg *config.Config) error {
	adminURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(adminURL)))
	defer sqldb.Close()

	ctx := context.Background()
	if err := sqldb.PingContext(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	var exists bool
	row := sqldb.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DBName,
	)
	if err := row.Scan(&exists); err != nil {
		return fmt.Errorf("check database existence: %w", err)
	}

	if !exists {
		_, err := sqldb.ExecContext(ctx,
			fmt.Sprintf("CREATE DATABASE %q", cfg.DBName),
		)
		if err != nil {
			return fmt.Errorf("create database %q: %w", cfg.DBName, err)
		}
	}

	return nil
}
