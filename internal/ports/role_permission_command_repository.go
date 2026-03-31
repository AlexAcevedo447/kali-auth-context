package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type IAssignPermissionToRoleCommandRepository interface {
	Assign(ctx context.Context, relation *identity.RolePermission) error
}

type IRemovePermissionFromRoleCommandRepository interface {
	Remove(ctx context.Context, tenantId identity.TenantId, roleId identity.RoleId, permissionId identity.PermissionId) error
}
