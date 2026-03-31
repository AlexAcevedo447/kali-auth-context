package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type IGetUserRolesQueryRepository interface {
	GetByUser(ctx context.Context, tenantId identity.TenantId, userId identity.UserId) ([]*identity.UserRole, error)
}
