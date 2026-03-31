package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type IGetTenantByIdQueryRepository interface {
	GetById(ctx context.Context, tenantId identity.TenantId) (*identity.Tenant, error)
}

type IListTenantsQueryRepository interface {
	List(ctx context.Context) ([]*identity.Tenant, error)
}

type IGetTenantByNameQueryRepository interface {
	GetByName(ctx context.Context, name string) (*identity.Tenant, error)
}
