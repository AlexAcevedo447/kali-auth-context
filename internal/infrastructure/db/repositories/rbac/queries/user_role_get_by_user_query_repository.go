package queries

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type GetUserRolesQueryRepository struct {
	db *bun.DB
}

var _ ports.IGetUserRolesQueryRepository = (*GetUserRolesQueryRepository)(nil)

func NewGetUserRolesQueryRepository(db *bun.DB) *GetUserRolesQueryRepository {
	return &GetUserRolesQueryRepository{db: db}
}

func (r *GetUserRolesQueryRepository) GetByUser(ctx context.Context, tenantId identity.TenantId, userId identity.UserId) ([]*identity.UserRole, error) {
	var list []models.UserRoleModel
	if err := r.db.NewSelect().
		Model(&list).
		Where("tenant_id = ? AND user_id = ?", string(tenantId), string(userId)).
		Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]*identity.UserRole, 0, len(list))
	for i := range list {
		result = append(result, mappers.ToDomainUserRole(&list[i]))
	}

	return result, nil
}
