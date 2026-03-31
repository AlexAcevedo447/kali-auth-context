package db

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetTenantByNameQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetTenantByNameQueryRepository = (*GetTenantByNameQueryRepository)(nil)

func NewGetTenantByNameQueryRepository(db *bun.DB) *GetTenantByNameQueryRepository {
	return &GetTenantByNameQueryRepository{db: db}
}

func (r *GetTenantByNameQueryRepository) GetByName(ctx context.Context, name string) (*identity.Tenant, error) {
	var model models.TenantModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("name = ?", name).
		Scan(ctx); err != nil {
		return nil, err
	}

	return mappers.ToDomainTenant(&model), nil
}
