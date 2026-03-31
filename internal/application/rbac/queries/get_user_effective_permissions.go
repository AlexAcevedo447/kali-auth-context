package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetUserEffectivePermissionsQuery struct {
	userRolesRepo       ports.IGetUserRolesQueryRepository
	rolePermissionsRepo ports.IGetRolePermissionsQueryRepository
	permissionRepo      ports.IGetPermissionByIdQueryRepository
}

func NewGetUserEffectivePermissionsQuery(
	userRolesRepo ports.IGetUserRolesQueryRepository,
	rolePermissionsRepo ports.IGetRolePermissionsQueryRepository,
	permissionRepo ports.IGetPermissionByIdQueryRepository,
) *GetUserEffectivePermissionsQuery {
	return &GetUserEffectivePermissionsQuery{
		userRolesRepo:       userRolesRepo,
		rolePermissionsRepo: rolePermissionsRepo,
		permissionRepo:      permissionRepo,
	}
}

func (q *GetUserEffectivePermissionsQuery) Handle(ctx context.Context, tenantId identity.TenantId, userId identity.UserId) ([]*identity.Permission, error) {
	if tenantId == "" {
		return nil, identity.ErrTenantRequired
	}
	if userId == "" {
		return nil, identity.ErrUserIdRequired
	}

	userRoles, err := q.userRolesRepo.GetByUser(ctx, tenantId, userId)
	if err != nil {
		return nil, err
	}

	permissionById := make(map[identity.PermissionId]*identity.Permission)

	for _, userRole := range userRoles {
		rolePermissions, err := q.rolePermissionsRepo.GetByRole(ctx, tenantId, userRole.RoleId)
		if err != nil {
			return nil, err
		}

		for _, rolePermission := range rolePermissions {
			if _, exists := permissionById[rolePermission.PermissionId]; exists {
				continue
			}

			perm, err := q.permissionRepo.GetById(ctx, rolePermission.PermissionId)
			if err != nil {
				return nil, err
			}

			permissionById[rolePermission.PermissionId] = perm
		}
	}

	result := make([]*identity.Permission, 0, len(permissionById))
	for _, p := range permissionById {
		result = append(result, p)
	}

	return result, nil
}
