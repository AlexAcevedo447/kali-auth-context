package db

import (
	"context"
	"fmt"
	"log"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
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

	// Código legible de empresa: columna, valores para tenants existentes e índice único.
	if _, err := m.db.ExecContext(ctx, "ALTER TABLE tenants ADD COLUMN IF NOT EXISTS code VARCHAR"); err != nil {
		return err
	}
	if err := m.backfillTenantCodes(ctx); err != nil {
		return err
	}
	if _, err := m.db.ExecContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS idx_tenants_code_unique ON tenants (code)"); err != nil {
		return err
	}

	// El login resuelve el tenant a partir del correo.
	if _, err := m.db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_users_email ON user_models (email)"); err != nil {
		return err
	}

	// Integridad RBAC (incremento 2): índices únicos solo si los datos existentes
	// lo permiten. Nunca se borran ni fusionan datos aquí; si hay duplicados, el
	// índice queda pendiente (la aplicación igual impide nuevos duplicados).
	if _, err := m.EnsureRBACIntegrityIndexes(ctx); err != nil {
		return err
	}

	return nil
}

// RBACIntegrityIndex describe una restricción de unicidad RBAC y cómo detectar
// los datos que impiden crearla.
type RBACIntegrityIndex struct {
	Name       string
	Create     string
	Duplicates string
}

var RBACIntegrityIndexes = []RBACIntegrityIndex{
	{
		Name:       "idx_roles_tenant_name_unique",
		Create:     "CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_tenant_name_unique ON roles (tenant_id, lower(trim(name)))",
		Duplicates: "SELECT count(*) FROM (SELECT 1 FROM roles GROUP BY tenant_id, lower(trim(name)) HAVING count(*) > 1) d",
	},
	{
		Name:       "idx_permissions_tenant_resource_action_unique",
		Create:     "CREATE UNIQUE INDEX IF NOT EXISTS idx_permissions_tenant_resource_action_unique ON permissions (tenant_id, lower(trim(resource)), lower(trim(action)))",
		Duplicates: "SELECT count(*) FROM (SELECT 1 FROM permissions GROUP BY tenant_id, lower(trim(resource)), lower(trim(action)) HAVING count(*) > 1) d",
	},
	{
		Name:       "idx_user_roles_unique",
		Create:     "CREATE UNIQUE INDEX IF NOT EXISTS idx_user_roles_unique ON user_roles (tenant_id, user_id, role_id)",
		Duplicates: "SELECT count(*) FROM (SELECT 1 FROM user_roles GROUP BY tenant_id, user_id, role_id HAVING count(*) > 1) d",
	},
	{
		Name:       "idx_role_permissions_unique",
		Create:     "CREATE UNIQUE INDEX IF NOT EXISTS idx_role_permissions_unique ON role_permissions (tenant_id, role_id, permission_id)",
		Duplicates: "SELECT count(*) FROM (SELECT 1 FROM role_permissions GROUP BY tenant_id, role_id, permission_id HAVING count(*) > 1) d",
	},
}

// EnsureRBACIntegrityIndexes crea cada índice único si no hay duplicados y
// devuelve los nombres de los que quedaron pendientes por datos existentes.
func (m *Migrator) EnsureRBACIntegrityIndexes(ctx context.Context) ([]string, error) {
	var pending []string
	for _, idx := range RBACIntegrityIndexes {
		var dups int
		if err := m.db.NewRaw(idx.Duplicates).Scan(ctx, &dups); err != nil {
			return nil, err
		}
		if dups > 0 {
			log.Printf("WARN rbac integrity: %s pendiente, %d grupos duplicados (ver plan de migración en MATRIZ_PERMISOS.md)", idx.Name, dups)
			pending = append(pending, idx.Name)
			continue
		}
		if _, err := m.db.ExecContext(ctx, idx.Create); err != nil {
			return nil, err
		}
	}
	return pending, nil
}

// backfillTenantCodes asigna un código legible (derivado del nombre) a los
// tenants que aún no lo tienen, evitando colisiones con un sufijo numérico.
func (m *Migrator) backfillTenantCodes(ctx context.Context) error {
	var pending []models.TenantModel
	if err := m.db.NewSelect().Model(&pending).Where("code IS NULL OR code = ''").OrderExpr("id").Scan(ctx); err != nil {
		return err
	}
	for _, t := range pending {
		base := identity.SlugifyTenantCode(t.Name)
		code := base
		for n := 2; ; n++ {
			exists, err := m.db.NewSelect().Model((*models.TenantModel)(nil)).Where("code = ?", code).Exists(ctx)
			if err != nil {
				return err
			}
			if !exists {
				break
			}
			code = fmt.Sprintf("%s-%d", base, n)
		}
		if _, err := m.db.NewUpdate().Model((*models.TenantModel)(nil)).Set("code = ?", code).Where("id = ?", t.Id).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}
