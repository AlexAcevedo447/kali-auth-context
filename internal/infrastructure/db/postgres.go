package db

import (
	"context"
	"fmt"

	"kali-auth-context/internal/infrastructure/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	databaseUrl := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)

	poolCfg, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(context.Background(), poolCfg)
}
