package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type RemoveRoleFromUserCommandRepository struct {
	db *bun.DB
}

var _ ports.IRemoveRoleFromUserCommandRepository = (*RemoveRoleFromUserCommandRepository)(nil)

func NewRemoveRoleFromUserCommandRepository(db *bun.DB) *RemoveRoleFromUserCommandRepository {
	return &RemoveRoleFromUserCommandRepository{db: db}
}

func (r *RemoveRoleFromUserCommandRepository) Remove(ctx context.Context, tenantId identity.TenantId, userId identity.UserId, roleId identity.RoleId) error {
	_, err := r.db.NewDelete().
		Model((*models.UserRoleModel)(nil)).
		Where("tenant_id = ? AND user_id = ? AND role_id = ?", string(tenantId), string(userId), string(roleId)).
		Exec(ctx)
	return err
}
