package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type ICreateRoleCommandRepository interface {
	Create(ctx context.Context, role *identity.Role) error
}

type IUpdateRoleCommandRepository interface {
	Update(ctx context.Context, role *identity.Role) error
}

type IDeleteRoleCommandRepository interface {
	Delete(ctx context.Context, roleId identity.RoleId) error
}
