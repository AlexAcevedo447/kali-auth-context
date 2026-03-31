package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type IGetRoleByIdQueryRepository interface {
	GetById(ctx context.Context, roleId identity.RoleId) (*identity.Role, error)
}

type IListRolesQueryRepository interface {
	List(ctx context.Context, tenantId identity.TenantId) ([]*identity.Role, error)
}
