package ports

import "kali-auth-context/internal/domain/identity"

type IGetUserByIdQueryRepository interface {
	GetById(userId identity.UserId) (*identity.User, error)
}

type IGetUserByEmailQueryRepository interface {
	GetByEmail(email string) (*identity.User, error)
}

type IListUsersQueryRepository interface {
	List() ([]*identity.User, error)
}
