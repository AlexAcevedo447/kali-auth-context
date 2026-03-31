package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type IAssignRoleToUserCommandRepository interface {
	Assign(ctx context.Context, relation *identity.UserRole) error
}

type IRemoveRoleFromUserCommandRepository interface {
	Remove(ctx context.Context, tenantId identity.TenantId, userId identity.UserId, roleId identity.RoleId) error
}
