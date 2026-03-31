package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type IGetRolePermissionsQueryRepository interface {
	GetByRole(ctx context.Context, tenantId identity.TenantId, roleId identity.RoleId) ([]*identity.RolePermission, error)
}
