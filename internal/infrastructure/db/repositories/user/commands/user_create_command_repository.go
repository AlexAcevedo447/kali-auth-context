package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type CreateUserCommandRepository struct {
	db *bun.DB
}

var _ ports.ICreateUserCommandRepository = (*CreateUserCommandRepository)(nil)

func NewCreateUserCommandRepository(db *bun.DB) *CreateUserCommandRepository {
	return &CreateUserCommandRepository{db: db}
}

func (r *CreateUserCommandRepository) Create(ctx context.Context, user *identity.User) error {
	model := mappers.ToUserModel(user)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}
