package queries

import (
	"context"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetRoleByIdQuery struct {
	repo ports.IGetRoleByIdQueryRepository
}

func NewGetRoleByIdQuery(repo ports.IGetRoleByIdQueryRepository) *GetRoleByIdQuery {
	return &GetRoleByIdQuery{repo: repo}
}

func (q *GetRoleByIdQuery) Handle(ctx context.Context, roleId identity.RoleId) (*identity.Role, error) {
	if roleId == "" {
		return nil, identity.ErrRoleIdRequired
	}
	return q.repo.GetById(ctx, roleId)
}
