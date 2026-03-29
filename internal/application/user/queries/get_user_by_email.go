package queries

import (
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type GetUserByEmailQuery struct {
	query  ports.IGetUserByEmailQueryRepository
}

func NewGetUserByEmailQuery(query ports.IGetUserByEmailQueryRepository) *GetUserByEmailQuery {
	return &GetUserByEmailQuery{
		query: query,
	}
}

func (q *GetUserByEmailQuery) Handle(email string) (*identity.User, error) {
	return q.query.GetByEmail(email)
}