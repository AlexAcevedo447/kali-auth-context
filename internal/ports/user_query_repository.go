package ports

import "kali-auth-context/internal/domain/identity"

type IGetUserByIdQueryRepository interface {
	GetById(tenantId identity.TenantId, userId identity.UserId) (*identity.User, error)
}

type IGetUserByEmailQueryRepository interface {
	GetByEmail(tenantId identity.TenantId, email string) (*identity.User, error)
}

type IListUsersQueryRepository interface {
	List(tenantId identity.TenantId) ([]*identity.User, error)
}
