package queries

import (
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetUserByIdQuery struct {
	query  ports.IGetUserByIdQueryRepository
}

func NewGetUserByIdQuery(query ports.IGetUserByIdQueryRepository) *GetUserByIdQuery {
	return &GetUserByIdQuery{
		query: query,
	}
}

func (q *GetUserByIdQuery) Handle(tenantId identity.TenantId, id identity.UserId) (*identity.User, error) {
	return q.query.GetById(tenantId, identity.UserId(id))
}