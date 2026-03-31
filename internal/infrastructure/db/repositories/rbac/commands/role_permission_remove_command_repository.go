package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/models"
	"kali-auth-context/internal/ports"
)

type RemovePermissionFromRoleCommandRepository struct {
	db *bun.DB
}

var _ ports.IRemovePermissionFromRoleCommandRepository = (*RemovePermissionFromRoleCommandRepository)(nil)

func NewRemovePermissionFromRoleCommandRepository(db *bun.DB) *RemovePermissionFromRoleCommandRepository {
	return &RemovePermissionFromRoleCommandRepository{db: db}
}

func (r *RemovePermissionFromRoleCommandRepository) Remove(ctx context.Context, tenantId identity.TenantId, roleId identity.RoleId, permissionId identity.PermissionId) error {
	_, err := r.db.NewDelete().
		Model((*models.RolePermissionModel)(nil)).
		Where("tenant_id = ? AND role_id = ? AND permission_id = ?", string(tenantId), string(roleId), string(permissionId)).
		Exec(ctx)
	return err
}
