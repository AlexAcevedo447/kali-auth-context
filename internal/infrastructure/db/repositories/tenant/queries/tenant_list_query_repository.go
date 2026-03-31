package db

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type ListTenantsQueryRepository struct {
	db *bun.DB
}

var _ ports.IListTenantsQueryRepository = (*ListTenantsQueryRepository)(nil)

func NewListTenantsQueryRepository(db *bun.DB) *ListTenantsQueryRepository {
	return &ListTenantsQueryRepository{db: db}
}

func (r *ListTenantsQueryRepository) List(ctx context.Context) ([]*identity.Tenant, error) {
	var modelsList []models.TenantModel
	if err := r.db.NewSelect().Model(&modelsList).Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]*identity.Tenant, 0, len(modelsList))
	for i := range modelsList {
		result = append(result, mappers.ToDomainTenant(&modelsList[i]))
	}

	return result, nil
}
