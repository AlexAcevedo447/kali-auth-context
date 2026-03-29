package db

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"kali-auth-context/internal/infrastructure/config"
)

func NewPool(cfg *config.Config) (*bun.DB, error) {
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(databaseUrl),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	return db, nil
}
