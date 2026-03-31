package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type UpdatePermissionCommandRepository struct {
	db *bun.DB
}

var _ ports.IUpdatePermissionCommandRepository = (*UpdatePermissionCommandRepository)(nil)

func NewUpdatePermissionCommandRepository(db *bun.DB) *UpdatePermissionCommandRepository {
	return &UpdatePermissionCommandRepository{db: db}
}

func (r *UpdatePermissionCommandRepository) Update(ctx context.Context, permission *identity.Permission) error {
	model := mappers.ToPermissionModel(permission)
	_, err := r.db.NewUpdate().
		Model(model).
		Column("resource", "action").
		Where("id = ?", model.Id).
		Exec(ctx)
	return err
}
