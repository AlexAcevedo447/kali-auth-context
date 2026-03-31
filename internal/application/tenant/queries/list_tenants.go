package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type ListTenantsQuery struct {
	repo ports.IListTenantsQueryRepository
}

func NewListTenantsQuery(repo ports.IListTenantsQueryRepository) *ListTenantsQuery {
	return &ListTenantsQuery{repo: repo}
}

func (q *ListTenantsQuery) Handle(ctx context.Context) ([]*identity.Tenant, error) {
	return q.repo.List(ctx)
}
