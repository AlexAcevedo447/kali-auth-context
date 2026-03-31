package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetUserRolesQuery struct {
	repo ports.IGetUserRolesQueryRepository
}

func NewGetUserRolesQuery(repo ports.IGetUserRolesQueryRepository) *GetUserRolesQuery {
	return &GetUserRolesQuery{repo: repo}
}

func (q *GetUserRolesQuery) Handle(ctx context.Context, tenantId identity.TenantId, userId identity.UserId) ([]*identity.UserRole, error) {
	if tenantId == "" {
		return nil, identity.ErrTenantRequired
	}
	if userId == "" {
		return nil, identity.ErrUserIdRequired
	}
	return q.repo.GetByUser(ctx, tenantId, userId)
}
