package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetTenantByIdQuery struct {
	repo ports.IGetTenantByIdQueryRepository
}

func NewGetTenantByIdQuery(repo ports.IGetTenantByIdQueryRepository) *GetTenantByIdQuery {
	return &GetTenantByIdQuery{repo: repo}
}

func (q *GetTenantByIdQuery) Handle(ctx context.Context, tenantId identity.TenantId) (*identity.Tenant, error) {
	if tenantId == "" {
		return nil, identity.ErrTenantRequired
	}

	return q.repo.GetById(ctx, tenantId)
}
