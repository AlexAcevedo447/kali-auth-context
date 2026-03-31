package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetRolePermissionsQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetRolePermissionsQueryRepository = (*GetRolePermissionsQueryRepository)(nil)

func NewGetRolePermissionsQueryRepository(db *bun.DB) *GetRolePermissionsQueryRepository {
	return &GetRolePermissionsQueryRepository{db: db}
}

func (r *GetRolePermissionsQueryRepository) GetByRole(ctx context.Context, tenantId identity.TenantId, roleId identity.RoleId) ([]*identity.RolePermission, error) {
	var list []models.RolePermissionModel
	if err := r.db.NewSelect().
		Model(&list).
		Where("tenant_id = ? AND role_id = ?", string(tenantId), string(roleId)).
		Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]*identity.RolePermission, 0, len(list))
	for i := range list {
		result = append(result, mappers.ToDomainRolePermission(&list[i]))
	}

	return result, nil
}
