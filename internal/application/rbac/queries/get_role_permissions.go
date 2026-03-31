package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetRolePermissionsQuery struct {
	repo ports.IGetRolePermissionsQueryRepository
}

func NewGetRolePermissionsQuery(repo ports.IGetRolePermissionsQueryRepository) *GetRolePermissionsQuery {
	return &GetRolePermissionsQuery{repo: repo}
}

func (q *GetRolePermissionsQuery) Handle(ctx context.Context, tenantId identity.TenantId, roleId identity.RoleId) ([]*identity.RolePermission, error) {
	if tenantId == "" {
		return nil, identity.ErrTenantRequired
	}
	if roleId == "" {
		return nil, identity.ErrRoleIdRequired
	}
	return q.repo.GetByRole(ctx, tenantId, roleId)
}
