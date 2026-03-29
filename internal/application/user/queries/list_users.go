package queries

import (
	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/ports"
)

type ListUsersQuery struct {
	query  ports.IListUsersQueryRepository
}

func NewListUsersQuery(query ports.IListUsersQueryRepository) *ListUsersQuery {
	return &ListUsersQuery{
		query: query,
	}
}

func (q *ListUsersQuery) Handle() ([]*identity.User, error) {
	return q.query.List()
}