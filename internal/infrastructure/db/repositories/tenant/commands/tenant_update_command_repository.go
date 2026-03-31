package db

import (
	"context"

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
	_, err := r.db.NewUpdate().
		Model(model).
		Column("name", "status").
		Where("id = ?", model.Id).
		Exec(ctx)
	return err
}
