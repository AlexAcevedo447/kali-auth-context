package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetTenantByNameQuery struct {
	repo ports.IGetTenantByNameQueryRepository
}

func NewGetTenantByNameQuery(repo ports.IGetTenantByNameQueryRepository) *GetTenantByNameQuery {
	return &GetTenantByNameQuery{repo: repo}
}

func (q *GetTenantByNameQuery) Handle(ctx context.Context, name string) (*identity.Tenant, error) {
	return q.repo.GetByName(ctx, name)
}
