package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type ListPermissionsQueryRepository struct {
	db *bun.DB
}

var _ ports.IListPermissionsQueryRepository = (*ListPermissionsQueryRepository)(nil)

func NewListPermissionsQueryRepository(db *bun.DB) *ListPermissionsQueryRepository {
	return &ListPermissionsQueryRepository{db: db}
}

func (r *ListPermissionsQueryRepository) List(ctx context.Context, tenantId identity.TenantId) ([]*identity.Permission, error) {
	var list []models.PermissionModel
	if err := r.db.NewSelect().
		Model(&list).
		Where("tenant_id = ?", string(tenantId)).
		Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]*identity.Permission, 0, len(list))
	for i := range list {
		result = append(result, mappers.ToDomainPermission(&list[i]))
	}

	return result, nil
}
