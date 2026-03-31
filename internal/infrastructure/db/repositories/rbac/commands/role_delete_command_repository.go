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
	_, err := r.db.NewDelete().
		Model((*models.RoleModel)(nil)).
		Where("id = ?", string(roleId)).
		Exec(ctx)
	return err
}
