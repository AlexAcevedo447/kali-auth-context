package commands

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type DeleteUserCommandRepository struct {
	db *bun.DB
}

var _ ports.IDeleteUserCommandRepository = (*DeleteUserCommandRepository)(nil)

func NewDeleteUserCommandRepository(db *bun.DB) *DeleteUserCommandRepository {
	return &DeleteUserCommandRepository{db: db}
}

func (r *DeleteUserCommandRepository) Delete(ctx context.Context, tenantId identity.TenantId, userId identity.UserId) error {
	// Lo único que referencia a un usuario son sus asignaciones de rol
	// (user_roles). Sin cascada: si aún tiene roles, 409 y se quitan antes con
	// las reglas RBAC. Así no quedan asignaciones huérfanas.
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		hasRoles, err := tx.NewSelect().TableExpr("user_roles").
			Where("tenant_id = ? AND user_id = ?", tenantId, userId).
			Exists(ctx)
		if err != nil {
			return err
		}
		if hasRoles {
			return identity.ErrUserHasRoles
		}
		res, err := tx.NewDelete().
			Model((*models.UserModel)(nil)).
			Where("tenant_id = ? AND id = ?", tenantId, userId).
			Exec(ctx)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}
