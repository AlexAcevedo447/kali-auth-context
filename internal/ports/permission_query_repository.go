package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type IGetPermissionByIdQueryRepository interface {
	GetById(ctx context.Context, permissionId identity.PermissionId) (*identity.Permission, error)
}

type IListPermissionsQueryRepository interface {
	List(ctx context.Context, tenantId identity.TenantId) ([]*identity.Permission, error)
}
