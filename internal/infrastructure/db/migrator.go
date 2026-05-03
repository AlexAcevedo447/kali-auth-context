package db

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/infrastructure/db/models"
)

// Migrator is an infrastructure component responsible for ensuring the database
// schema is up to date. It belongs to the infrastructure layer, not the entry
// point, so that main.go stays clean and the composition root is properly
// separated from bootstrap concerns.
type Migrator struct {
	db *bun.DB
}

func NewMigrator(db *bun.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate(ctx context.Context) error {
	_, err := m.db.NewCreateTable().
		Model((*models.IdempotencyModel)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}
