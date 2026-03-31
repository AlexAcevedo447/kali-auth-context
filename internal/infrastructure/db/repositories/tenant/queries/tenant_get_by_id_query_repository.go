package db

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetTenantByIdQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetTenantByIdQueryRepository = (*GetTenantByIdQueryRepository)(nil)

func NewGetTenantByIdQueryRepository(db *bun.DB) *GetTenantByIdQueryRepository {
	return &GetTenantByIdQueryRepository{db: db}
}

func (r *GetTenantByIdQueryRepository) GetById(ctx context.Context, tenantId identity.TenantId) (*identity.Tenant, error) {
	var model models.TenantModel
	if err := r.db.NewSelect().
		Model(&model).
		Where("id = ?", string(tenantId)).
		Scan(ctx); err != nil {
		return nil, err
	}

	return mappers.ToDomainTenant(&model), nil
}
