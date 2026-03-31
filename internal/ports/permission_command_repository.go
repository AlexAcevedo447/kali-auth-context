package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type ICreatePermissionCommandRepository interface {
	Create(ctx context.Context, permission *identity.Permission) error
}

type IUpdatePermissionCommandRepository interface {
	Update(ctx context.Context, permission *identity.Permission) error
}

type IDeletePermissionCommandRepository interface {
	Delete(ctx context.Context, permissionId identity.PermissionId) error
}
