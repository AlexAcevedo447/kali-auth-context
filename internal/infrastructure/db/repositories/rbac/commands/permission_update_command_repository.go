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
	// Editar un permiso asignado cambiaría los accesos de otros usuarios.
	assigned, err := exists(ctx, r.db, "SELECT 1 FROM role_permissions WHERE permission_id = ?", string(permission.Id))
	if err != nil {
		return err
	}
	if assigned {
		return identity.ErrPermissionAssigned
	}
	taken, err := permissionTaken(ctx, r.db, string(permission.TenantId), permission.Resource, permission.Action, string(permission.Id))
	if err != nil {
		return err
	}
	if taken {
		return identity.ErrPermissionTaken
	}
	model := mappers.ToPermissionModel(permission)
	_, err = r.db.NewUpdate().
		Model(model).
		Column("resource", "action").
		Where("id = ?", model.Id).
		Exec(ctx)
	return err
}
