package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type DeleteRoleCommandRepository struct {
	db *bun.DB
}

var _ ports.IDeleteRoleCommandRepository = (*DeleteRoleCommandRepository)(nil)

func NewDeleteRoleCommandRepository(db *bun.DB) *DeleteRoleCommandRepository {
	return &DeleteRoleCommandRepository{db: db}
}

func (r *DeleteRoleCommandRepository) Delete(ctx context.Context, roleId identity.RoleId) error {
	// Bloquea si hay usuarios con el rol; si no, borra el rol y sus vínculos
	// rol↔permiso en una transacción (no afecta los accesos de nadie).
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		assigned, err := exists(ctx, tx, "SELECT 1 FROM user_roles WHERE role_id = ?", string(roleId))
		if err != nil {
			return err
		}
		if assigned {
			return identity.ErrRoleAssigned
		}
		if _, err := tx.NewDelete().TableExpr("role_permissions").Where("role_id = ?", string(roleId)).Exec(ctx); err != nil {
			return err
		}
		_, err = tx.NewDelete().
			Model((*models.RoleModel)(nil)).
			Where("id = ?", string(roleId)).
			Exec(ctx)
		return err
	})
}
