package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type CreateRoleCommandRepository struct {
	db *bun.DB
}

var _ ports.ICreateRoleCommandRepository = (*CreateRoleCommandRepository)(nil)

func NewCreateRoleCommandRepository(db *bun.DB) *CreateRoleCommandRepository {
	return &CreateRoleCommandRepository{db: db}
}

func (r *CreateRoleCommandRepository) Create(ctx context.Context, role *identity.Role) error {
	model := mappers.ToRoleModel(role)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}
