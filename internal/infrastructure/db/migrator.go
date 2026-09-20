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

	// Unicidad de email por tenant (no global): dos usuarios del mismo tenant
	// no pueden compartir email, pero el mismo email sí puede repetirse entre
	// tenants distintos. Defensa definitiva contra condiciones de carrera,
	// más allá de cualquier verificación previa a nivel de aplicación.
	if _, err := m.db.NewCreateIndex().
		Model((*models.UserModel)(nil)).
		Index("idx_users_tenant_id_email_unique").
		Column("tenant_id", "email").
		Unique().
		IfNotExists().
		Exec(ctx); err != nil {
		return err
	}

	return nil
}
