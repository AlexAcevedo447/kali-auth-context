package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetPermissionByIdQuery struct {
	repo ports.IGetPermissionByIdQueryRepository
}

func NewGetPermissionByIdQuery(repo ports.IGetPermissionByIdQueryRepository) *GetPermissionByIdQuery {
	return &GetPermissionByIdQuery{repo: repo}
}

func (q *GetPermissionByIdQuery) Handle(ctx context.Context, permissionId identity.PermissionId) (*identity.Permission, error) {
	if permissionId == "" {
		return nil, identity.ErrPermissionIdRequired
	}
	return q.repo.GetById(ctx, permissionId)
}
