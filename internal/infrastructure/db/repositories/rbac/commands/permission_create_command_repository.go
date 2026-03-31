package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type CreatePermissionCommandRepository struct {
	db *bun.DB
}

var _ ports.ICreatePermissionCommandRepository = (*CreatePermissionCommandRepository)(nil)

func NewCreatePermissionCommandRepository(db *bun.DB) *CreatePermissionCommandRepository {
	return &CreatePermissionCommandRepository{db: db}
}

func (r *CreatePermissionCommandRepository) Create(ctx context.Context, permission *identity.Permission) error {
	model := mappers.ToPermissionModel(permission)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}
