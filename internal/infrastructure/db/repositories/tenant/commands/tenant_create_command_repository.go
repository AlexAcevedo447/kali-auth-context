package db

import (
	"context"

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
	return err
}
