package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type ListRolesQuery struct {
	repo ports.IListRolesQueryRepository
}

func NewListRolesQuery(repo ports.IListRolesQueryRepository) *ListRolesQuery {
	return &ListRolesQuery{repo: repo}
}

func (q *ListRolesQuery) Handle(ctx context.Context, tenantId identity.TenantId) ([]*identity.Role, error) {
	if tenantId == "" {
		return nil, identity.ErrTenantRequired
	}
	return q.repo.List(ctx, tenantId)
}
