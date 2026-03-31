package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type ListRolesQueryRepository struct {
	db *bun.DB
}

var _ ports.IListRolesQueryRepository = (*ListRolesQueryRepository)(nil)

func NewListRolesQueryRepository(db *bun.DB) *ListRolesQueryRepository {
	return &ListRolesQueryRepository{db: db}
}

func (r *ListRolesQueryRepository) List(ctx context.Context, tenantId identity.TenantId) ([]*identity.Role, error) {
	var list []models.RoleModel
	if err := r.db.NewSelect().
		Model(&list).
		Where("tenant_id = ?", string(tenantId)).
		Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]*identity.Role, 0, len(list))
	for i := range list {
		result = append(result, mappers.ToDomainRole(&list[i]))
	}

	return result, nil
}
