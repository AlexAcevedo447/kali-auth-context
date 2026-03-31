package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type DeletePermissionCommandRepository struct {
	db *bun.DB
}

var _ ports.IDeletePermissionCommandRepository = (*DeletePermissionCommandRepository)(nil)

func NewDeletePermissionCommandRepository(db *bun.DB) *DeletePermissionCommandRepository {
	return &DeletePermissionCommandRepository{db: db}
}

func (r *DeletePermissionCommandRepository) Delete(ctx context.Context, permissionId identity.PermissionId) error {
	_, err := r.db.NewDelete().
		Model((*models.PermissionModel)(nil)).
		Where("id = ?", string(permissionId)).
		Exec(ctx)
	return err
}
