package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type ICreateTenantCommandRepository interface {
	Create(ctx context.Context, tenant *identity.Tenant) error
}

type IUpdateTenantCommandRepository interface {
	Update(ctx context.Context, tenant *identity.Tenant) error
}

type IActivateTenantCommandRepository interface {
	Activate(ctx context.Context, tenantId identity.TenantId) error
}

type ISuspendTenantCommandRepository interface {
	Suspend(ctx context.Context, tenantId identity.TenantId) error
}
