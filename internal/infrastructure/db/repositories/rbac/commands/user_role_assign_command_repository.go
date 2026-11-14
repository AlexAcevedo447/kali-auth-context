package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type AssignRoleToUserCommandRepository struct {
	db *bun.DB
}

var _ ports.IAssignRoleToUserCommandRepository = (*AssignRoleToUserCommandRepository)(nil)

func NewAssignRoleToUserCommandRepository(db *bun.DB) *AssignRoleToUserCommandRepository {
	return &AssignRoleToUserCommandRepository{db: db}
}

func (r *AssignRoleToUserCommandRepository) Assign(ctx context.Context, relation *identity.UserRole) error {
	// Idempotente: si ya está asignado no se duplica.
	already, err := exists(ctx, r.db,
		"SELECT 1 FROM user_roles WHERE tenant_id = ? AND user_id = ? AND role_id = ?",
		string(relation.TenantId), string(relation.UserId), string(relation.RoleId))
	if err != nil || already {
		return err
	}
	model := mappers.ToUserRoleModel(relation)
	_, err = r.db.NewInsert().Model(model).Exec(ctx)
	return err
}
