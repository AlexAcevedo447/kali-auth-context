package ports

import (
	"context"

	"kali-auth-context/internal/domain/identity"
)

type ICreateUserCommandRepository interface {
	Create(ctx context.Context, user *identity.User) error
}

type IUpdateUserCommandRepository interface {
	Update(ctx context.Context, user *identity.User) error
}

type IDeleteUserCommandRepository interface {
	Delete(ctx context.Context, userId identity.UserId) error
}
