package commands

import (
	"context"

	"github.com/uptrace/bun"
)

// exists ejecuta SELECT EXISTS(...) sobre una consulta con sus argumentos.
func exists(ctx context.Context, db bun.IDB, query string, args ...any) (bool, error) {
	var found bool
	err := db.NewRaw("SELECT EXISTS ("+query+")", args...).Scan(ctx, &found)
	return found, err
}

// Unicidad de nombre de rol por empresa, sin distinguir mayúsculas ni espacios
// exteriores. excludeID permite editar el propio rol sin chocar consigo mismo.
func roleNameTaken(ctx context.Context, db bun.IDB, tenantID, name, excludeID string) (bool, error) {
	return exists(ctx, db,
		"SELECT 1 FROM roles WHERE tenant_id = ? AND lower(trim(name)) = lower(trim(?)) AND id <> ?",
		tenantID, name, excludeID)
}

// Unicidad de permiso por (empresa, recurso, acción) normalizados.
func permissionTaken(ctx context.Context, db bun.IDB, tenantID, resource, action, excludeID string) (bool, error) {
	return exists(ctx, db,
		"SELECT 1 FROM permissions WHERE tenant_id = ? AND lower(trim(resource)) = lower(trim(?)) AND lower(trim(action)) = lower(trim(?)) AND id <> ?",
		tenantID, resource, action, excludeID)
}
