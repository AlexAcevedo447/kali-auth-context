package commands

import (
	"context"

	"github.com/uptrace/bun"

	"kali-auth-context/internal/domain/identity"
	"kali-auth-context/internal/infrastructure/db/mappers"
	"kali-auth-context/internal/ports"
)

type AssignPermissionToRoleCommandRepository struct {
	db *bun.DB
}

var _ ports.IAssignPermissionToRoleCommandRepository = (*AssignPermissionToRoleCommandRepository)(nil)

func NewAssignPermissionToRoleCommandRepository(db *bun.DB) *AssignPermissionToRoleCommandRepository {
	return &AssignPermissionToRoleCommandRepository{db: db}
}

func (r *AssignPermissionToRoleCommandRepository) Assign(ctx context.Context, relation *identity.RolePermission) error {
	model := mappers.ToRolePermissionModel(relation)
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}
