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
	assigned, err := exists(ctx, r.db, "SELECT 1 FROM role_permissions WHERE permission_id = ?", string(permissionId))
	if err != nil {
		return err
	}
	if assigned {
		return identity.ErrPermissionAssigned
	}
	_, err = r.db.NewDelete().
		Model((*models.PermissionModel)(nil)).
		Where("id = ?", string(permissionId)).
		Exec(ctx)
	return err
}
