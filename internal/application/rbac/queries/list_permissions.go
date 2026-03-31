package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type ListPermissionsQuery struct {
	repo ports.IListPermissionsQueryRepository
}

func NewListPermissionsQuery(repo ports.IListPermissionsQueryRepository) *ListPermissionsQuery {
	return &ListPermissionsQuery{repo: repo}
}

func (q *ListPermissionsQuery) Handle(ctx context.Context, tenantId identity.TenantId) ([]*identity.Permission, error) {
	if tenantId == "" {
		return nil, identity.ErrTenantRequired
	}
	return q.repo.List(ctx, tenantId)
}
