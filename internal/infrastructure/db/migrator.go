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
	tables := []interface{}{
		(*models.IdempotencyModel)(nil),
		(*models.TenantModel)(nil),
		(*models.UserModel)(nil),
		(*models.RoleModel)(nil),
		(*models.PermissionModel)(nil),
		(*models.UserRoleModel)(nil),
		(*models.RolePermissionModel)(nil),
	}

	for _, table := range tables {
		if _, err := m.db.NewCreateTable().Model(table).IfNotExists().Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}
