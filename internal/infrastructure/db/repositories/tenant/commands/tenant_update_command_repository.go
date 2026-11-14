package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type UpdateTenantCommandRepository struct {
	db *bun.DB
}

var _ ports.IUpdateTenantCommandRepository = (*UpdateTenantCommandRepository)(nil)

func NewUpdateTenantCommandRepository(db *bun.DB) *UpdateTenantCommandRepository {
	return &UpdateTenantCommandRepository{db: db}
}

func (r *UpdateTenantCommandRepository) Update(ctx context.Context, tenant *identity.Tenant) error {
	model := mappers.ToTenantModel(tenant)
	columns := []string{"name", "status"}
	if model.Code != "" {
		columns = append(columns, "code")
	}
	res, err := r.db.NewUpdate().
		Model(model).
		Column(columns...).
		Where("id = ?", model.Id).
		Exec(ctx)
	if err != nil {
		return translateTenantCodeConflict(err)
	}
	return requireAffectedRow(res)
}

// requireAffectedRow devuelve sql.ErrNoRows (→ 404) si la sentencia no modificó
// ninguna fila, es decir, si la empresa no existe.
func requireAffectedRow(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
