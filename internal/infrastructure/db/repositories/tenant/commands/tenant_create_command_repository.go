package db

import (
	"context"
	"errors"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type CreateTenantCommandRepository struct {
	db *bun.DB
}

var _ ports.ICreateTenantCommandRepository = (*CreateTenantCommandRepository)(nil)

func NewCreateTenantCommandRepository(db *bun.DB) *CreateTenantCommandRepository {
	return &CreateTenantCommandRepository{db: db}
}

func (r *CreateTenantCommandRepository) Create(ctx context.Context, tenant *identity.Tenant) error {
	model := mappers.ToTenantModel(tenant)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return translateTenantCodeConflict(err)
}

// pgUniqueViolation es el SQLSTATE de Postgres para una violación de unicidad.
const pgUniqueViolation = "23505"

// translateTenantCodeConflict traduce la violación del índice único del código de
// empresa (idx_tenants_code_unique) al error de dominio ErrTenantCodeTaken (409).
// Se usa en crear y actualizar. Se compara por la interfaz Field (que implementa
// pgdriver.Error) para poder probarlo sin base de datos.
func translateTenantCodeConflict(err error) error {
	var pgErr interface{ Field(k byte) string }
	if errors.As(err, &pgErr) && pgErr.Field('C') == pgUniqueViolation {
		return identity.ErrTenantCodeTaken
	}
	return err
}
